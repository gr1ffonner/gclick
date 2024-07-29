package db

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	UUID              uuid.UUID `json:"uuid" validate:"required,uuid"`
	MaterialType      string    `json:"material_type" validate:"required,oneof='статья' 'видеоролик' 'презентация'"`
	PublicationStatus string    `json:"publication_status" validate:"required,oneof='архивный' 'активный'"`
	Title             string    `json:"title" validate:"required,max=255"`
	Content           string    `json:"content" validate:"omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
