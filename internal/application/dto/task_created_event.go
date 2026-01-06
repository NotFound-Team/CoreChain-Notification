package dto

import (
	"time"

	"github.com/corechain/notification-service/internal/domain/models"
)

type TaskCreatedEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      models.Task            `json:"data"`
	Metadata  TaskEventMetadata      `json:"metadata"`
}

type TaskEventMetadata struct {
	AssignedToUser AssignedUserInfo `json:"assignedToUser"`
}

type AssignedUserInfo struct {
	ID       string `json:"_id"`
	Email    string `json:"email"`
	FCMToken string `json:"fcmToken"`
	Name     string `json:"name"`
}

type MessageCreatedEvent struct {
	EventType string         `json:"event_type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      models.Message `json:"data"`
	Metadata  interface{}    `json:"metadata"`
}

type IncomingCallEvent struct {
	EventType string       `json:"event_type"`
	Timestamp time.Time    `json:"timestamp"`
	Data      models.Call  `json:"data"`
	Metadata  interface{}  `json:"metadata"`
}
