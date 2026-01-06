package interfaces

import (
	"context"

	"github.com/corechain/notification-service/internal/domain/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	Update(ctx context.Context, notification *models.Notification) error
	GetByID(ctx context.Context, id string) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error)
	GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error)
	UpdateStatus(ctx context.Context, id string, status string, errorMsg string) error
}

type FCMClient interface {
	SendNotification(ctx context.Context, token string, title string, body string, data map[string]string) error
	SendBatchNotifications(ctx context.Context, notifications []FCMMessage) error
}

type FCMMessage struct {
	Token string
	Title string
	Body  string
	Data  map[string]string
}

type KafkaConsumer interface {
	Start(ctx context.Context) error
	Stop() error
	RegisterHandler(topic string, handler MessageHandler) error
}

type MessageHandler func(ctx context.Context, message []byte) error
