package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/corechain/notification-service/internal/application/services"
	"github.com/corechain/notification-service/internal/config"
	"github.com/corechain/notification-service/internal/delivery/kafka"
	httpDelivery "github.com/corechain/notification-service/internal/delivery/http"
	"github.com/corechain/notification-service/internal/delivery/http/handlers"
	"github.com/corechain/notification-service/internal/infrastructure/fcm"
	kafkaInfra "github.com/corechain/notification-service/internal/infrastructure/kafka"
	"github.com/corechain/notification-service/internal/infrastructure/repository/postgres"
	"github.com/corechain/notification-service/internal/utils/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	if err := cfg.Validate(); err != nil {
		panic("Invalid configuration: " + err.Error())
	}

	if err := logger.Init(cfg.Logger.Level); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting Notification Service",
		zap.String("env", cfg.Server.Env),
		zap.String("log_level", cfg.Logger.Level),
		zap.Int("port", cfg.Server.Port),
	)

	ctx := context.Background()

	logger.Info("Connecting to PostgreSQL...")
	repository, err := postgres.NewNotificationRepository(cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer repository.Close()
	logger.Info("Successfully connected to PostgreSQL")

	logger.Info("Initializing FCM client...")
	fcmClient, err := fcm.NewClient(ctx, cfg.FCM.CredentialsPath, cfg.FCM.ProjectID)
	if err != nil {
		logger.Fatal("Failed to initialize FCM client", zap.Error(err))
	}
	logger.Info("Successfully initialized FCM client",
		zap.String("project_id", cfg.FCM.ProjectID),
	)

	notificationService := services.NewNotificationService(repository, fcmClient)
	taskNotificationService := services.NewTaskNotificationService(notificationService)

	// Initialize HTTP server
	logger.Info("Initializing HTTP server...")
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	httpServer := httpDelivery.NewServer(httpDelivery.ServerConfig{
		Port:                cfg.Server.Port,
		NotificationHandler: notificationHandler,
	})

	logger.Info("Initializing Kafka consumer...")
	kafkaConsumer := kafkaInfra.NewConsumer(kafkaInfra.ConsumerConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topics:  []string{cfg.Kafka.Topics.TaskCreated},
	})

	taskHandler := kafka.NewTaskHandler(taskNotificationService)
	if err := kafkaConsumer.RegisterHandler(cfg.Kafka.Topics.TaskCreated, taskHandler.HandleTaskCreated); err != nil {
		logger.Fatal("Failed to register task.created handler", zap.Error(err))
	}

	logger.Info("Registered Kafka handlers",
		zap.String("task_created_topic", cfg.Kafka.Topics.TaskCreated),
	)

	// Start HTTP server in a goroutine
	go func() {
		if err := httpServer.Start(cfg.Server.Port); err != nil {
			logger.Error("HTTP server error", zap.Error(err))
		}
	}()
	logger.Info("HTTP server started", zap.Int("port", cfg.Server.Port))

	// Start Kafka consumer
	logger.Info("Starting Kafka consumer...")
	if err := kafkaConsumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start Kafka consumer", zap.Error(err))
	}

	logger.Info("Notification Service started successfully",
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers),
		zap.String("consumer_group", cfg.Kafka.GroupID),
		zap.Int("http_port", cfg.Server.Port),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received, stopping service...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop Kafka consumer
	if err := kafkaConsumer.Stop(); err != nil {
		logger.Error("Error stopping Kafka consumer", zap.Error(err))
	}

	// Stop HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error stopping HTTP server", zap.Error(err))
	}

	logger.Info("Notification Service stopped")
}
