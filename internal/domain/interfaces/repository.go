package interfaces

import (
	"context"

	"github.com/corechain/notification-service/internal/domain/models"
)

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	// Create saves a new notification record
	Create(ctx context.Context, notification *models.Notification) error
	
	// Update updates an existing notification
	Update(ctx context.Context, notification *models.Notification) error
	
	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id string) (*models.Notification, error)
	
	// GetByUserID retrieves notifications for a specific user
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error)
	
	// GetPendingNotifications retrieves all pending notifications
	GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error)
	
	// UpdateStatus updates the status of a notification
	UpdateStatus(ctx context.Context, id string, status string, errorMsg string) error
}

// FCMClient defines the interface for Firebase Cloud Messaging operations
type FCMClient interface {
	// SendNotification sends a single notification via FCM
	SendNotification(ctx context.Context, token string, title string, body string, data map[string]string) error
	
	// SendBatchNotifications sends multiple notifications via FCM
	SendBatchNotifications(ctx context.Context, notifications []FCMMessage) error
}

// FCMMessage represents a message to be sent via FCM
type FCMMessage struct {
	Token string
	Title string
	Body  string
	Data  map[string]string
}

// KafkaConsumer defines the interface for Kafka message consumption
type KafkaConsumer interface {
	// Start begins consuming messages from Kafka
	Start(ctx context.Context) error
	
	// Stop gracefully stops consuming messages
	Stop() error
	
	// RegisterHandler registers a handler for a specific topic
	RegisterHandler(topic string, handler MessageHandler) error
}

// MessageHandler is a function type for handling Kafka messages
type MessageHandler func(ctx context.Context, message []byte) error
