package dto

type FCMNotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
type FCMDataPayload struct {
	Type        string `json:"type"`
	TaskID      string `json:"task_id,omitempty"`
	Priority    string `json:"priority,omitempty"`
	ClickAction string `json:"click_action"`
}
type FCMAndroidConfig struct {
	Priority     string                        `json:"priority"`
	Notification FCMAndroidNotificationConfig  `json:"notification"`
}

type FCMAndroidNotificationConfig struct {
	Sound     string `json:"sound"`
	ChannelID string `json:"channel_id"`
}

type FCMAPNSPayload struct {
	Payload FCMAPSConfig `json:"payload"`
}

type FCMAPSConfig struct {
	APS FCMAPSContent `json:"aps"`
}

type FCMAPSContent struct {
	Sound string `json:"sound"`
	Badge int    `json:"badge"`
}

type FCMMessage struct {
	Token        string                  `json:"token"`
	Notification FCMNotificationPayload  `json:"notification"`
	Data         map[string]string       `json:"data"`
	Android      FCMAndroidConfig        `json:"android"`
	APNS         FCMAPNSPayload          `json:"apns"`
}
