package services

import (
	"context"

	"github.com/corechain/notification-service/internal/application/dto"
	"github.com/corechain/notification-service/internal/domain/models"
	"github.com/corechain/notification-service/internal/infrastructure/fcm"
	"github.com/corechain/notification-service/pkg/constants"
)

// TaskNotificationService handles task-specific notification logic
type TaskNotificationService struct {
	notificationService *NotificationService
}

// NewTaskNotificationService creates a new task notification service
func NewTaskNotificationService(notificationService *NotificationService) *TaskNotificationService {
	return &TaskNotificationService{
		notificationService: notificationService,
	}
}

// ProcessTaskCreatedEvent processes a task.created event from Kafka
func (s *TaskNotificationService) ProcessTaskCreatedEvent(ctx context.Context, event *dto.TaskCreatedEvent) error {
	// Build notification template
	template := fcm.BuildTaskCreatedNotification(
		event.Data.Title,
		event.Data.CreatedBy.Email,
		event.Data.DueDate,
		event.Data.Priority,
	)

	// Prepare notification data
	data := map[string]interface{}{
		"type":        "task_created",
		"task_id":     event.Data.ID,
		"project_id":  event.Data.ProjectID,
		"priority":    event.Data.Priority,
		"click_action": "OPEN_TASK_DETAIL",
	}

	notification := &models.Notification{
		NotificationType: constants.NotificationTypeTaskCreated,
		UserID:           event.Metadata.AssignedToUser.ID,
		FCMToken:         event.Metadata.AssignedToUser.FCMToken,
		Title:            template.Title,
		Body:             template.Body,
		Data:             data,
		TaskID:           event.Data.ID,
		ProjectID:        event.Data.ProjectID,
		Priority:         event.Data.Priority,
	}

	// Send notification
	return s.notificationService.CreateAndSendNotification(ctx, notification)
}

// ProcessTaskUpdatedEvent processes a task.updated event from Kafka (future)
func (s *TaskNotificationService) ProcessTaskUpdatedEvent(ctx context.Context, event *dto.TaskCreatedEvent) error {
	// Similar to task created, but with different template
	template := fcm.BuildTaskUpdatedNotification(
		event.Data.Title,
		event.Data.UpdatedBy.Email,
	)

	data := map[string]interface{}{
		"type":        "task_updated",
		"task_id":     event.Data.ID,
		"project_id":  event.Data.ProjectID,
		"click_action": "OPEN_TASK_DETAIL",
	}

	notification := &models.Notification{
		NotificationType: constants.NotificationTypeTaskUpdated,
		UserID:           event.Metadata.AssignedToUser.ID,
		FCMToken:         event.Metadata.AssignedToUser.FCMToken,
		Title:            template.Title,
		Body:             template.Body,
		Data:             data,
		TaskID:           event.Data.ID,
		ProjectID:        event.Data.ProjectID,
		Priority:         event.Data.Priority,
	}

	return s.notificationService.CreateAndSendNotification(ctx, notification)
}
