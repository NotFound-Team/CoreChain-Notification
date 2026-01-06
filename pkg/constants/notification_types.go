package constants

type NotificationType string

const (
	NotificationTypeTaskCreated NotificationType = "task_created"
	
	NotificationTypeTaskUpdated NotificationType = "task_updated"
	
	NotificationTypeNewMessage NotificationType = "new_message"
	
	NotificationTypeIncomingCall NotificationType = "incoming_call"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	
	StatusSent NotificationStatus = "sent"
	
	StatusFailed NotificationStatus = "failed"
)

const (
	PriorityHigh   = 1
	PriorityMedium = 2
	PriorityLow    = 3
)

const (
	TaskStatusPending    = 0
	TaskStatusInProgress = 1
	TaskStatusCompleted  = 2
)
