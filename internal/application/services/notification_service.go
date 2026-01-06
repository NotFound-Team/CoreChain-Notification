package services

import (
	"context"
	"fmt"
	"time"

	"github.com/corechain/notification-service/internal/domain/interfaces"
	"github.com/corechain/notification-service/internal/domain/models"
	"github.com/corechain/notification-service/internal/utils/logger"
	"github.com/corechain/notification-service/pkg/constants"
	"go.uber.org/zap"
)

type NotificationService struct {
	repository interfaces.NotificationRepository
	fcmClient  interfaces.FCMClient
}

func NewNotificationService(repo interfaces.NotificationRepository, fcmClient interfaces.FCMClient) *NotificationService {
	return &NotificationService{
		repository: repo,
		fcmClient:  fcmClient,
	}
}

func (s *NotificationService) CreateAndSendNotification(ctx context.Context, notification *models.Notification) error {
	notification.CreatedAt = time.Now()
	notification.Status = constants.StatusPending

	if err := s.repository.Create(ctx, notification); err != nil {
		logger.Error("Failed to create notification in database",
			zap.Error(err),
			zap.String("user_id", notification.UserID),
		)
		return err
	}

	logger.Info("Created notification record",
		zap.String("id", notification.ID),
		zap.String("type", string(notification.NotificationType)),
		zap.String("user_id", notification.UserID),
	)

	dataMap := make(map[string]string)
	for k, v := range notification.Data {
		dataMap[k] = fmt.Sprintf("%v", v)
	}

	if err := s.fcmClient.SendNotification(ctx, notification.FCMToken, notification.Title, notification.Body, dataMap); err != nil {
		// Update status to failed
		updateErr := s.repository.UpdateStatus(ctx, notification.ID, string(constants.StatusFailed), err.Error())
		if updateErr != nil {
			logger.Error("Failed to update notification status to failed",
				zap.Error(updateErr),
				zap.String("notification_id", notification.ID),
			)
		}

		logger.Error("Failed to send FCM notification",
			zap.Error(err),
			zap.String("notification_id", notification.ID),
			zap.String("user_id", notification.UserID),
		)
		return err
	}

	if err := s.repository.UpdateStatus(ctx, notification.ID, string(constants.StatusSent), ""); err != nil {
		logger.Error("Failed to update notification status to sent",
			zap.Error(err),
			zap.String("notification_id", notification.ID),
		)
		return err
	}

	logger.Info("Successfully sent notification",
		zap.String("id", notification.ID),
		zap.String("user_id", notification.UserID),
		zap.String("type", string(notification.NotificationType)),
	)

	return nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error) {
	return s.repository.GetByUserID(ctx, userID, limit, offset)
}

func (s *NotificationService) RetryFailedNotifications(ctx context.Context, maxRetries int) error {
	notifications, err := s.repository.GetPendingNotifications(ctx, 100)
	if err != nil {
		return err
	}

	for _, notif := range notifications {
		if notif.RetryCount >= maxRetries {
			logger.Warn("Max retries exceeded for notification",
				zap.String("notification_id", notif.ID),
				zap.Int("retry_count", notif.RetryCount),
			)
			continue
		}

		if err := s.CreateAndSendNotification(ctx, notif); err != nil {
			logger.Error("Failed to retry notification",
				zap.Error(err),
				zap.String("notification_id", notif.ID),
			)
		}
	}

	return nil
}
