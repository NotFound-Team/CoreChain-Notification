package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/corechain/notification-service/internal/domain/interfaces"
	"github.com/corechain/notification-service/internal/utils/errors"
	"google.golang.org/api/option"
)

// Client implements the FCM client interface
type Client struct {
	messagingClient *messaging.Client
}

// NewClient creates a new FCM client
func NewClient(ctx context.Context, credentialsPath string) (*Client, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, errors.NewFCMError("failed to initialize Firebase app", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, errors.NewFCMError("failed to initialize FCM messaging client", err)
	}

	return &Client{
		messagingClient: messagingClient,
	}, nil
}

// SendNotification sends a single notification via FCM
func (c *Client) SendNotification(ctx context.Context, token string, title string, body string, data map[string]string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Sound:     "default",
				ChannelID: "task_notifications",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: intPtr(1),
				},
			},
		},
	}

	response, err := c.messagingClient.Send(ctx, message)
	if err != nil {
		return errors.NewFCMError(fmt.Sprintf("failed to send FCM notification to token %s", token), err)
	}

	// Log successful send
	_ = response // Response contains message ID

	return nil
}

// SendBatchNotifications sends multiple notifications via FCM
func (c *Client) SendBatchNotifications(ctx context.Context, notifications []interfaces.FCMMessage) error {
	if len(notifications) == 0 {
		return nil
	}

	messages := make([]*messaging.Message, len(notifications))
	for i, notif := range notifications {
		messages[i] = &messaging.Message{
			Token: notif.Token,
			Notification: &messaging.Notification{
				Title: notif.Title,
				Body:  notif.Body,
			},
			Data: notif.Data,
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Sound:     "default",
					ChannelID: "task_notifications",
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Sound: "default",
						Badge: intPtr(1),
					},
				},
			},
		}
	}

	batchResponse, err := c.messagingClient.SendAll(ctx, messages)
	if err != nil {
		return errors.NewFCMError("failed to send batch notifications", err)
	}

	// Check for failures
	if batchResponse.FailureCount > 0 {
		return errors.NewFCMError(
			fmt.Sprintf("failed to send %d out of %d notifications", batchResponse.FailureCount, len(notifications)),
			nil,
		)
	}

	return nil
}

func intPtr(i int) *int {
	return &i
}
