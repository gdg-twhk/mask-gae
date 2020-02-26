package service

import (
	"context"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/feedback/model"
	"github.com/cage1016/mask/internal/pkg/errors"
)

const QueryDatefmt = "2006_0102"

var (
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrInvalidQueryParams indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrInvalidQueryParams = errors.New("invalid query params")
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(FeedbacksvcService) FeedbacksvcService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type FeedbacksvcService interface {
	// [method=get,expose=true,router=api/feedback/options]
	Options(ctx context.Context) (items []model.Option, err error)
	// [method=get,expose=true,router=api/feedback/pharmacies/:pharmacie_id]
	PharmacyFeedBacks(ctx context.Context, PharmacyID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error)
	// [method=get,expose=true,router=api/feedback/users/:user_id]
	UserFeedBacks(ctx context.Context, userID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error)
	// [method=post,expose=true,router=api/feedback]
	InsertFeedBack(ctx context.Context, userID, pharmacyID, optionID, description string, Longitude, Latitude float64) (id string, err error)
}

// the concrete implementation of service interface
type stubFeedbacksvcService struct {
	logger  log.Logger
	repo    model.FeedbackRepository
	idpNano NanoIdentityProvider
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.FeedbackRepository, idpNano NanoIdentityProvider, logger log.Logger) (s FeedbacksvcService) {
	var svc FeedbacksvcService
	{
		svc = &stubFeedbacksvcService{repo: repo, idpNano: idpNano, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Options
func (fe *stubFeedbacksvcService) Options(ctx context.Context) (items []model.Option, err error) {
	return fe.repo.ListOption(ctx)
}

// Implement the business logic of PharmacyFeedBacks
func (fe *stubFeedbacksvcService) PharmacyFeedBacks(ctx context.Context, PharmacyID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	return fe.repo.RetrieveByPharmacyID(ctx, PharmacyID, date, offset, limit)
}

// Implement the business logic of UserFeedBacks
func (fe *stubFeedbacksvcService) UserFeedBacks(ctx context.Context, userID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	return fe.repo.RetrieveByUserID(ctx, userID, date, offset, limit)
}

// Implement the business logic of InsertFeedBack
func (fe *stubFeedbacksvcService) InsertFeedBack(ctx context.Context, userID, pharmacyID, optionID, description string, Longitude, Latitude float64) (id string, err error) {
	nid, _ := fe.idpNano.ID()
	return fe.repo.Insert(ctx, model.Feedback{
		ID:          nid,
		UserID:      userID,
		PharmacyID:  pharmacyID,
		OptionID:    optionID,
		Description: description,
		Longitude:   Longitude,
		Latitude:    Latitude,
	})
}
