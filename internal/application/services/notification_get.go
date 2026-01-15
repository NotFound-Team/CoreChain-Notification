package services

import (
	"context"

	"github.com/corechain/notification-service/internal/domain/models"
)

// GetNotificationByID retrieves a notification by its ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, id string) (*models.Notification, error) {
	return s.repository.GetByID(ctx, id)
}
