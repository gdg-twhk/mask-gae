package model

import (
	"context"
	"encoding/json"
	"time"
)

type Option struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Feedback struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	PharmacyID  string    `json:"storeId" db:"pharmacy_id"`
	OptionID    string    `json:"optionId" db:"option_id"`
	Description string    `json:"description" db:"description"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func (p *Feedback) MarshalJSON() ([]byte, error) {
	type Alias Feedback
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05-0700"),
	})
}

type FeedbackItemPage struct {
	PageMetadata
	Items []Feedback `json:"items"`
}

// FeedbackRepository specifies an account persistence API.
type FeedbackRepository interface {
	// Insert persists the user account. A non-nil error is returned to indicate
	Insert(context.Context, Feedback) (string, error)

	// RetrieveByUserID retrieves user by its unique identifier (i.e. email).
	RetrieveByUserID(context.Context, string, string, uint64, uint64) (FeedbackItemPage, error)

	// RetrieveByPharmacyID retrieves user by its unique identifier (i.e. email, provider).
	RetrieveByPharmacyID(context.Context, string, string, uint64, uint64) (FeedbackItemPage, error)

	// ListOption
	ListOption(context.Context) ([]Option, error)
}
