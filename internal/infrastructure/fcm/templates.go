package fcm

import (
	"fmt"
	"time"
)

// Template represents a notification template
type Template struct {
	Title string
	Body  string
}

// BuildTaskCreatedNotification creates a notification for task creation
func BuildTaskCreatedNotification(taskTitle string, createdByName string, dueDate *time.Time, priority int) Template {
	priorityText := getPriorityText(priority)
	
	title := "New Task Assigned"
	body := fmt.Sprintf("%s - %s", taskTitle, priorityText)
	
	if dueDate != nil {
		dueDateStr := dueDate.Format("Jan 02")
		body = fmt.Sprintf("%s - Due: %s", taskTitle, dueDateStr)
	}
	
	if createdByName != "" {
		body = fmt.Sprintf("%s assigned you a task: %s", createdByName, taskTitle)
	}

	return Template{
		Title: title,
		Body:  body,
	}
}

// BuildTaskUpdatedNotification creates a notification for task updates
func BuildTaskUpdatedNotification(taskTitle string, updatedByName string) Template {
	title := "Task Updated"
	body := fmt.Sprintf("%s has been updated", taskTitle)
	
	if updatedByName != "" {
		body = fmt.Sprintf("%s updated: %s", updatedByName, taskTitle)
	}

	return Template{
		Title: title,
		Body:  body,
	}
}

// BuildNewMessageNotification creates a notification for new messages (future)
func BuildNewMessageNotification(senderName string, messagePreview string) Template {
	title := fmt.Sprintf("New message from %s", senderName)
	body := messagePreview

	if len(body) > 100 {
		body = body[:97] + "..."
	}

	return Template{
		Title: title,
		Body:  body,
	}
}

// BuildIncomingCallNotification creates a notification for incoming calls (future)
func BuildIncomingCallNotification(callerName string, callType string) Template {
	title := "Incoming Call"
	body := fmt.Sprintf("%s is calling (%s)", callerName, callType)

	return Template{
		Title: title,
		Body:  body,
	}
}

func getPriorityText(priority int) string {
	switch priority {
	case 1:
		return "High Priority"
	case 2:
		return "Medium Priority"
	case 3:
		return "Low Priority"
	default:
		return ""
	}
}
