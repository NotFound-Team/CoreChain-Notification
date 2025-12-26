package dto

// FCMNotificationPayload represents the notification part of FCM message
type FCMNotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// FCMDataPayload represents the data part of FCM message
type FCMDataPayload struct {
	Type        string `json:"type"`
	TaskID      string `json:"task_id,omitempty"`
	Priority    string `json:"priority,omitempty"`
	ClickAction string `json:"click_action"`
}

// FCMAndroidConfig represents Android-specific configuration
type FCMAndroidConfig struct {
	Priority     string                        `json:"priority"`
	Notification FCMAndroidNotificationConfig  `json:"notification"`
}

// FCMAndroidNotificationConfig represents Android notification settings
type FCMAndroidNotificationConfig struct {
	Sound     string `json:"sound"`
	ChannelID string `json:"channel_id"`
}

// FCMAPNSPayload represents Apple Push Notification Service configuration
type FCMAPNSPayload struct {
	Payload FCMAPSConfig `json:"payload"`
}

// FCMAPSConfig represents the APS configuration
type FCMAPSConfig struct {
	APS FCMAPSContent `json:"aps"`
}

// FCMAPSContent represents the APS content
type FCMAPSContent struct {
	Sound string `json:"sound"`
	Badge int    `json:"badge"`
}

// FCMMessage represents the complete FCM message structure
type FCMMessage struct {
	Token        string                  `json:"token"`
	Notification FCMNotificationPayload  `json:"notification"`
	Data         map[string]string       `json:"data"`
	Android      FCMAndroidConfig        `json:"android"`
	APNS         FCMAPNSPayload          `json:"apns"`
}
