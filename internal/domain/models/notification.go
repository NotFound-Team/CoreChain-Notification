package models

import (
	"time"

	"github.com/corechain/notification-service/pkg/constants"
)

type NotificationType = constants.NotificationType
type NotificationStatus = constants.NotificationStatus
type Notification struct {
	ID               string                        `json:"id"`
	NotificationType NotificationType              `json:"notification_type"`
	UserID           string                        `json:"user_id"`
	FCMToken         string                        `json:"fcm_token"`
	Title            string                        `json:"title"`
	Body             string                        `json:"body"`
	Data             map[string]interface{}        `json:"data"`
	Status           NotificationStatus            `json:"status"`
	ErrorMessage     string                        `json:"error_message,omitempty"`
	CreatedAt        time.Time                     `json:"created_at"`
	SentAt           *time.Time                    `json:"sent_at,omitempty"`
	RetryCount       int                           `json:"retry_count"`
	TaskID           string                        `json:"task_id,omitempty"`
	ProjectID        string                        `json:"project_id,omitempty"`
	Priority         int                           `json:"priority,omitempty"`
}

type UserInfo struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Task struct {
	ID          string     `json:"_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Attachments []string   `json:"attachments"`
	CreatedBy   UserInfo   `json:"createdBy"`
	AssignedTo  string     `json:"assignedTo"`  // ObjectId reference
	ProjectID   string     `json:"projectId"`
	Priority    int        `json:"priority"`
	Status      int        `json:"status"`
	StartDate   *time.Time `json:"startDate"`
	DueDate     *time.Time `json:"dueDate"`
	IsDeleted   bool       `json:"isDeleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	UpdatedBy   *UserInfo  `json:"updatedBy,omitempty"`
	DeletedBy   *UserInfo  `json:"deletedBy,omitempty"`
}

type Message struct {
	ID        string    `json:"_id"`
	Content   string    `json:"content"`
	SenderID  string    `json:"senderId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Call struct {
	ID       string `json:"_id"`
	CallerID string `json:"callerId"`
	CallType string `json:"callType"`
}
