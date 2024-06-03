package db

type Article struct {
	EventID   int64  `json:"eventID" validate:"required"`
	EventType string `json:"eventType" validate:"required"`
	UserID    int64  `json:"userID" validate:"required"`
	EventTime string `json:"eventTime" validate:"required"`
	Payload   string `json:"payload" validate:"required"`
}
