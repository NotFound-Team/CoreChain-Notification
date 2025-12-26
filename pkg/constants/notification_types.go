package constants

// NotificationType defines the types of notifications
type NotificationType string

const (
	// NotificationTypeTaskCreated is sent when a new task is created
	NotificationTypeTaskCreated NotificationType = "task_created"
	
	// NotificationTypeTaskUpdated is sent when a task is updated
	NotificationTypeTaskUpdated NotificationType = "task_updated"
	
	// NotificationTypeNewMessage is sent when a new message arrives
	NotificationTypeNewMessage NotificationType = "new_message"
	
	// NotificationTypeIncomingCall is sent when there's an incoming call
	NotificationTypeIncomingCall NotificationType = "incoming_call"
)

// NotificationStatus represents the delivery status of a notification
type NotificationStatus string

const (
	// StatusPending means the notification is queued but not sent
	StatusPending NotificationStatus = "pending"
	
	// StatusSent means the notification was successfully sent
	StatusSent NotificationStatus = "sent"
	
	// StatusFailed means the notification failed to send
	StatusFailed NotificationStatus = "failed"
)

// Priority levels for tasks (matching NestJS schema)
const (
	PriorityHigh   = 1
	PriorityMedium = 2
	PriorityLow    = 3
)

// Task status values (matching NestJS schema)
const (
	TaskStatusPending    = 0
	TaskStatusInProgress = 1
	TaskStatusCompleted  = 2
)
