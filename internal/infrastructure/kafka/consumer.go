package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/corechain/notification-service/internal/domain/interfaces"
	"github.com/corechain/notification-service/internal/utils/errors"
	"github.com/corechain/notification-service/internal/utils/logger"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Consumer implements the Kafka consumer interface
type Consumer struct {
	readers  map[string]*kafkago.Reader
	handlers map[string]interfaces.MessageHandler
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
}

// ConsumerConfig holds Kafka consumer configuration
type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config ConsumerConfig) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	
	consumer := &Consumer{
		readers:  make(map[string]*kafkago.Reader),
		handlers: make(map[string]interfaces.MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}

	for _, topic := range config.Topics {
		reader := kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        config.Brokers,
			GroupID:        config.GroupID,
			Topic:          topic,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: 0, // Manual commit for reliability
		})
		consumer.readers[topic] = reader
	}

	return consumer
}

// RegisterHandler registers a handler for a specific topic
func (c *Consumer) RegisterHandler(topic string, handler interfaces.MessageHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.readers[topic]; !exists {
		return errors.NewKafkaError(fmt.Sprintf("no reader registered for topic: %s", topic), nil)
	}

	c.handlers[topic] = handler
	logger.Info("Registered handler for topic", zap.String("topic", topic))
	
	return nil
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for topic, reader := range c.readers {
		handler, exists := c.handlers[topic]
		if !exists {
			logger.Warn("No handler registered for topic, skipping", zap.String("topic", topic))
			continue
		}

		c.wg.Add(1)
		go c.consumeTopic(ctx, topic, reader, handler)
		logger.Info("Started consuming topic", zap.String("topic", topic))
	}

	logger.Info("Kafka consumer started successfully")
	return nil
}

// consumeTopic consumes messages from a specific topic
func (c *Consumer) consumeTopic(ctx context.Context, topic string, reader *kafkago.Reader, handler interfaces.MessageHandler) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping consumer for topic", zap.String("topic", topic))
			return
		default:
			message, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Error("Error fetching message",
					zap.String("topic", topic),
					zap.Error(err),
				)
				continue
			}

			logger.Debug("Received message",
				zap.String("topic", topic),
				zap.Int("partition", message.Partition),
				zap.Int64("offset", message.Offset),
			)

			if err := handler(ctx, message.Value); err != nil {
				logger.Error("Error processing message",
					zap.String("topic", topic),
					zap.Error(err),
					zap.ByteString("message", message.Value),
				)
				// Don't commit on error - ensures message is reprocessed
				continue
			}
			if err := reader.CommitMessages(ctx, message); err != nil {
				logger.Error("Error committing message",
					zap.String("topic", topic),
					zap.Error(err),
				)
			} else {
				logger.Debug("Successfully processed and committed message",
					zap.String("topic", topic),
					zap.Int64("offset", message.Offset),
				)
			}
		}
	}
}

// Stop gracefully stops consuming messages
func (c *Consumer) Stop() error {
	logger.Info("Stopping Kafka consumer...")
	
	c.cancel()
	c.wg.Wait()

	c.mu.Lock()
	defer c.mu.Unlock()
	
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			logger.Error("Error closing reader",
				zap.String("topic", topic),
				zap.Error(err),
			)
		}
	}

	logger.Info("Kafka consumer stopped")
	return nil
}
