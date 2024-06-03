package db

type Article struct {
	EventID int64  `json:"eventID"`
	UserID  int64  `json:"userID"`
	Payload string `json:"payload"`
	Validator
}

type Validator struct {
	EventType string `json:"eventType"`
	EventTime string `json:"eventTime"`
}
