package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/corechain/notification-service/internal/domain/models"
	"github.com/corechain/notification-service/internal/utils/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NotificationEntity represents the database entity
type NotificationEntity struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	NotificationType string    `gorm:"column:notification_type;type:varchar(50);not null"`
	UserID           string    `gorm:"column:user_id;type:varchar(100);not null;index"`
	FCMToken         string    `gorm:"column:fcm_token;type:text;not null"`
	Title            string    `gorm:"column:title;type:varchar(255);not null"`
	Body             string    `gorm:"column:body;type:text;not null"`
	Data             string    `gorm:"column:data;type:jsonb"`
	Status           string    `gorm:"column:status;type:varchar(20);not null;default:pending;index"`
	ErrorMessage     string    `gorm:"column:error_message;type:text"`
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:now();index:idx_notifications_created_at,sort:desc"`
	SentAt           *time.Time `gorm:"column:sent_at"`
	RetryCount       int       `gorm:"column:retry_count;default:0"`
	TaskID           string    `gorm:"column:task_id;type:varchar(100);index"`
	ProjectID        string    `gorm:"column:project_id;type:varchar(100)"`
	Priority         int       `gorm:"column:priority"`
}

// TableName specifies the table name
func (NotificationEntity) TableName() string {
	return "notifications"
}

// NotificationRepository implements the repository interface
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(dsn string) (*NotificationRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.NewDatabaseError("failed to connect to database", err)
	}

	return &NotificationRepository{db: db}, nil
}

// Create saves a new notification record
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	entity := r.toEntity(notification)
	
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return errors.NewDatabaseError("failed to create notification", err)
	}

	notification.ID = entity.ID
	return nil
}

// Update updates an existing notification
func (r *NotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	entity := r.toEntity(notification)
	
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return errors.NewDatabaseError("failed to update notification", err)
	}

	return nil
}

// GetByID retrieves a notification by ID
func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	var entity NotificationEntity
	
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.ErrCodeNotFound, "notification not found", err)
		}
		return nil, errors.NewDatabaseError("failed to get notification", err)
	}

	return r.toModel(&entity)
}

// GetByUserID retrieves notifications for a specific user
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, error) {
	var entities []NotificationEntity
	
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&entities).Error; err != nil {
		return nil, errors.NewDatabaseError("failed to get notifications by user ID", err)
	}

	var notifications []*models.Notification
	for _, entity := range entities {
		notification, err := r.toModel(&entity)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// GetPendingNotifications retrieves all pending notifications
func (r *NotificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]*models.Notification, error) {
	var entities []NotificationEntity
	
	query := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit)

	if err := query.Find(&entities).Error; err != nil {
		return nil, errors.NewDatabaseError("failed to get pending notifications", err)
	}

	var notifications []*models.Notification
	for _, entity := range entities {
		notification, err := r.toModel(&entity)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// UpdateStatus updates the status of a notification
func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "sent" {
		now := time.Now()
		updates["sent_at"] = now
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	if err := r.db.WithContext(ctx).Model(&NotificationEntity{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.NewDatabaseError("failed to update notification status", err)
	}

	return nil
}

// toEntity converts domain model to database entity
func (r *NotificationRepository) toEntity(notification *models.Notification) *NotificationEntity {
	entity := &NotificationEntity{
		ID:               notification.ID,
		NotificationType: string(notification.NotificationType),
		UserID:           notification.UserID,
		FCMToken:         notification.FCMToken,
		Title:            notification.Title,
		Body:             notification.Body,
		Status:           string(notification.Status),
		ErrorMessage:     notification.ErrorMessage,
		CreatedAt:        notification.CreatedAt,
		SentAt:           notification.SentAt,
		RetryCount:       notification.RetryCount,
		TaskID:           notification.TaskID,
		ProjectID:        notification.ProjectID,
		Priority:         notification.Priority,
	}

	if notification.Data != nil {
		dataJSON, _ := json.Marshal(notification.Data)
		entity.Data = string(dataJSON)
	}

	return entity
}

// toModel converts database entity to domain model
func (r *NotificationRepository) toModel(entity *NotificationEntity) (*models.Notification, error) {
	notification := &models.Notification{
		ID:               entity.ID,
		NotificationType: models.NotificationType(entity.NotificationType),
		UserID:           entity.UserID,
		FCMToken:         entity.FCMToken,
		Title:            entity.Title,
		Body:             entity.Body,
		Status:           models.NotificationStatus(entity.Status),
		ErrorMessage:     entity.ErrorMessage,
		CreatedAt:        entity.CreatedAt,
		SentAt:           entity.SentAt,
		RetryCount:       entity.RetryCount,
		TaskID:           entity.TaskID,
		ProjectID:        entity.ProjectID,
		Priority:         entity.Priority,
	}

	if entity.Data != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(entity.Data), &data); err != nil {
			return nil, errors.NewDatabaseError("failed to unmarshal notification data", err)
		}
		notification.Data = data
	}

	return notification, nil
}

// Close closes the database connection
func (r *NotificationRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
