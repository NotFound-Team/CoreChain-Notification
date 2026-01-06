package services

import (
	"context"

	"github.com/corechain/notification-service/internal/application/dto"
	"github.com/corechain/notification-service/internal/domain/models"
	"github.com/corechain/notification-service/internal/infrastructure/fcm"
	"github.com/corechain/notification-service/pkg/constants"
)

type TaskNotificationService struct {
	notificationService *NotificationService
}

func NewTaskNotificationService(notificationService *NotificationService) *TaskNotificationService {
	return &TaskNotificationService{
		notificationService: notificationService,
	}
}

func (s *TaskNotificationService) ProcessTaskCreatedEvent(ctx context.Context, event *dto.TaskCreatedEvent) error {
	template := fcm.BuildTaskCreatedNotification(
		event.Data.Title,
		event.Data.CreatedBy.Email,
		event.Data.DueDate,
		event.Data.Priority,
	)

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

	return s.notificationService.CreateAndSendNotification(ctx, notification)
}

func (s *TaskNotificationService) ProcessTaskUpdatedEvent(ctx context.Context, event *dto.TaskCreatedEvent) error {
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
