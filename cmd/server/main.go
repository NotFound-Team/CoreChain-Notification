package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/corechain/notification-service/internal/application/services"
	"github.com/corechain/notification-service/internal/config"
	"github.com/corechain/notification-service/internal/delivery/kafka"
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
	fcmClient, err := fcm.NewClient(ctx, cfg.FCM.CredentialsPath)
	if err != nil {
		logger.Fatal("Failed to initialize FCM client", zap.Error(err))
	}
	logger.Info("Successfully initialized FCM client",
		zap.String("project_id", cfg.FCM.ProjectID),
	)

	notificationService := services.NewNotificationService(repository, fcmClient)
	taskNotificationService := services.NewTaskNotificationService(notificationService)

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

	logger.Info("Starting Kafka consumer...")
	if err := kafkaConsumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start Kafka consumer", zap.Error(err))
	}

	logger.Info("Notification Service started successfully",
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers),
		zap.String("consumer_group", cfg.Kafka.GroupID),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received, stopping service...")

	if err := kafkaConsumer.Stop(); err != nil {
		logger.Error("Error stopping Kafka consumer", zap.Error(err))
	}

	logger.Info("Notification Service stopped")
}
