package db

type Article struct {
	EventType string `json:"eventType"`
	UserID    int    `json:"userID"`
	EventTime string `json:"eventTime"`
	Payload   string `json:"payload"`
}
