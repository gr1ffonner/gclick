package db

type Article struct {
	EventID   int64  `json:"eventID"`
	EventType string `json:"eventType"`
	UserID    int64  `json:"userID"`
	EventTime string `json:"eventTime"`
	Payload   string `json:"payload"`
}
