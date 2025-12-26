package dto

import (
	"time"

	"github.com/corechain/notification-service/internal/domain/models"
)

// TaskCreatedEvent represents the Kafka message for task.created event
type TaskCreatedEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      models.Task            `json:"data"`
	Metadata  TaskEventMetadata      `json:"metadata"`
}

// TaskEventMetadata contains additional metadata for the task event
type TaskEventMetadata struct {
	AssignedToUser AssignedUserInfo `json:"assignedToUser"`
}

// AssignedUserInfo contains user information including FCM token
type AssignedUserInfo struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	FCMToken string `json:"fcmToken"`
	Name     string `json:"name"`
}

// MessageCreatedEvent represents the Kafka message for message.new event (future)
type MessageCreatedEvent struct {
	EventType string         `json:"event_type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      models.Message `json:"data"`
	Metadata  interface{}    `json:"metadata"`
}

// IncomingCallEvent represents the Kafka message for call.incoming event (future)
type IncomingCallEvent struct {
	EventType string       `json:"event_type"`
	Timestamp time.Time    `json:"timestamp"`
	Data      models.Call  `json:"data"`
	Metadata  interface{}  `json:"metadata"`
}
