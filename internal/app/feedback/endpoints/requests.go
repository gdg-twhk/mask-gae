package endpoints

import (
	"time"

	"github.com/cage1016/mask/internal/app/feedback/service"
	"github.com/cage1016/mask/internal/pkg/errors"
)

const (
	maxLimitSize = 100

	customOptionID = "IRESxM58KC~dqg5XLCH~n"
)

type Request interface {
	validate() error
}

// OptionsRequest collects the request parameters for the Options method.
type OptionsRequest struct {
}

func (r OptionsRequest) validate() error {
	return nil // TBA
}

// PharmacyFeedBacksRequest collects the request parameters for the PharmacyFeedBacks method.
type PharmacyFeedBacksRequest struct {
	PharmacyID string `json:"pharmacyId"`
	Date       string `json:"date"`
	Limit      uint64 `json:"limit"`
	Offset     uint64 `json:"offset"`
}

func (r PharmacyFeedBacksRequest) validate() error {
	if r.PharmacyID == "" {
		return service.ErrMalformedEntity
	}

	_, err := time.Parse(service.QueryDatefmt, r.Date)
	if err != nil {
		return errors.Wrap(service.ErrMalformedEntity, err)
	}

	if r.Limit <= 0 || r.Limit > maxLimitSize {
		return errors.Wrap(service.ErrMalformedEntity, errors.New("limit must between 1 - 100"))
	}

	return nil
}

// UserFeedBacksRequest collects the request parameters for the UserFeedBacks method.
type UserFeedBacksRequest struct {
	UserID string `json:"user_id"`
	Date   string `json:"date"`
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

func (r UserFeedBacksRequest) validate() error {
	if r.UserID == "" {
		return service.ErrMalformedEntity
	}

	_, err := time.Parse(service.QueryDatefmt, r.Date)
	if err != nil {
		return errors.Wrap(service.ErrMalformedEntity, err)
	}

	if r.Limit <= 0 || r.Limit > maxLimitSize {
		return errors.Wrap(service.ErrMalformedEntity, errors.New("limit must between 1 - 100"))
	}

	return nil
}

// FeedBackRequest collects the request parameters for the InsertFeedBack method.
type FeedBackRequest struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	PharmacyID  string  `json:"pharmacyId"`
	OptionID    string  `json:"optionId"`
	Description string  `json:"description"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

func (r FeedBackRequest) validate() error {
	if r.UserID == "" || r.PharmacyID == "" || r.OptionID == "" {
		return errors.Wrap(service.ErrMalformedEntity, errors.New("userId or pharmacyId or optionId is empty"))
	}

	if r.OptionID == customOptionID && r.Description == "" {
		return errors.Wrap(service.ErrMalformedEntity, errors.New("description must have value when option is customize"))
	}

	return nil // TBA
}
