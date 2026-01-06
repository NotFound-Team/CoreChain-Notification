package kafka

import (
	"context"
	"encoding/json"

	"github.com/corechain/notification-service/internal/application/dto"
	"github.com/corechain/notification-service/internal/application/services"
	"github.com/corechain/notification-service/internal/utils/errors"
	"github.com/corechain/notification-service/internal/utils/logger"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskNotificationService *services.TaskNotificationService
}

func NewTaskHandler(taskNotificationService *services.TaskNotificationService) *TaskHandler {
	return &TaskHandler{
		taskNotificationService: taskNotificationService,
	}
}

func (h *TaskHandler) HandleTaskCreated(ctx context.Context, message []byte) error {
	logger.Debug("Processing task.created event", zap.Int("message_size", len(message)))

	var event dto.TaskCreatedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		logger.Error("Failed to unmarshal task.created event",
			zap.Error(err),
			zap.ByteString("message", message),
		)
		return errors.NewInvalidPayloadError("failed to unmarshal task.created event", err)
	}

	logger.Info("Received task.created event",
		zap.String("task_id", event.Data.ID),
		zap.String("assigned_to", event.Metadata.AssignedToUser.ID),
		zap.String("title", event.Data.Title),
	)

	if err := h.validateTaskCreatedEvent(&event); err != nil {
		logger.Error("Invalid task.created event", zap.Error(err))
		return err
	}

	if err := h.taskNotificationService.ProcessTaskCreatedEvent(ctx, &event); err != nil {
		logger.Error("Failed to process task.created event",
			zap.Error(err),
			zap.String("task_id", event.Data.ID),
		)
		return err
	}

	logger.Info("Successfully processed task.created event",
		zap.String("task_id", event.Data.ID),
	)

	return nil
}

func (h *TaskHandler) HandleTaskUpdated(ctx context.Context, message []byte) error {
	logger.Debug("Processing task.updated event", zap.Int("message_size", len(message)))

	var event dto.TaskCreatedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		logger.Error("Failed to unmarshal task.updated event",
			zap.Error(err),
			zap.ByteString("message", message),
		)
		return errors.NewInvalidPayloadError("failed to unmarshal task.updated event", err)
	}

	logger.Info("Received task.updated event",
		zap.String("task_id", event.Data.ID),
		zap.String("title", event.Data.Title),
	)

	if err := h.taskNotificationService.ProcessTaskUpdatedEvent(ctx, &event); err != nil {
		logger.Error("Failed to process task.updated event",
			zap.Error(err),
			zap.String("task_id", event.Data.ID),
		)
		return err
	}

	logger.Info("Successfully processed task.updated event",
		zap.String("task_id", event.Data.ID),
	)

	return nil
}

func (h *TaskHandler) validateTaskCreatedEvent(event *dto.TaskCreatedEvent) error {
	if event.Data.ID == "" {
		return errors.NewInvalidPayloadError("task ID is required", nil)
	}

	if event.Data.Title == "" {
		return errors.NewInvalidPayloadError("task title is required", nil)
	}

	if event.Metadata.AssignedToUser.ID == "" {
		return errors.NewInvalidPayloadError("assigned user ID is required", nil)
	}

	if event.Metadata.AssignedToUser.FCMToken == "" {
		return errors.NewInvalidPayloadError("FCM token is required", nil)
	}

	return nil
}
