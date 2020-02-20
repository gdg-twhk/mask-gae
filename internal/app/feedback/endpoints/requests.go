package endpoints

import (
	"time"

	"github.com/cage1016/mask/internal/app/feedback/model"
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

	if r.Limit == 0 || r.Limit > maxLimitSize {
		return service.ErrMalformedEntity
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

	if r.Limit == 0 || r.Limit > maxLimitSize {
		return service.ErrMalformedEntity
	}

	return nil
}

// FeedBackRequest collects the request parameters for the InsertFeedBack method.
type FeedBackRequest struct {
	Feedback model.Feedback `json:"feedback"`
}

func (r FeedBackRequest) validate() error {
	if r.Feedback.UserID == "" || r.Feedback.PharmacyID == "" || r.Feedback.OptionID == "" {
		return service.ErrMalformedEntity
	}

	if r.Feedback.OptionID == customOptionID && r.Feedback.Description == "" {
		return service.ErrMalformedEntity
	}

	return nil // TBA
}
