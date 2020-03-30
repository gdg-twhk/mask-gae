package service

import (
	"context"
	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/level"
	"github.com/go-kit/kit/log"
)

var (
	ErrInvalidTask     = errors.New("Bad Request - Invalid Task")
	ErrTaskCreatFailed = errors.New("task create failed")
	ErrMalformedEntity = errors.New("malformed entity specification")
)

const newLayout = "15:04"

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(PharmacyService) PharmacyService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type PharmacyService interface {
	// [method=post,expose=true,router=api/pharmacies]
	Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Pharmacy, err error)
	// [expose=false]
	TickerUpdate(ctx context.Context) (err error)
}

// the concrete implementation of service interface
type stubPharmacyService struct {
	logger              log.Logger
	repo                model.PharmacyRepository
	latestPharmacyTable string
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.PharmacyRepository, logger log.Logger) (s PharmacyService) {
	var svc PharmacyService
	{
		svc = &stubPharmacyService{repo: repo, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of TickerUpdate
func (as *stubPharmacyService) TickerUpdate(ctx context.Context) (err error) {
	return as._GetLatestPharmacyTableName(ctx)
}

func (as *stubPharmacyService) _GetLatestPharmacyTableName(ctx context.Context) (err error) {
	t, err := as.repo.GetLatestPharmacyTableName(ctx)
	if err != nil {
		level.Error(as.logger).Log("method", "as.repo.GetLatestPharmacyTableName", "err", err)
	}
	as.latestPharmacyTable = t

	as.logger.Log("latestPharmacyTable", t)
	return err
}

// Implement the business logic of Query
func (st *stubPharmacyService) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, _ float64, _ float64, swLng float64, swLat float64, _ float64, _ float64, max uint64) (items []model.Pharmacy, err error) {
	if st.latestPharmacyTable == "" {
		err := st._GetLatestPharmacyTableName(ctx)
		if err != nil {
			return []model.Pharmacy{}, err
		}
	}

	return st.repo.Query(ctx, st.latestPharmacyTable, centerLng, centerLat, swLng, neLng, swLat, neLat, max)
}
