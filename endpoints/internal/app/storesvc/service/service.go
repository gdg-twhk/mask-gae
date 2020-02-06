package service

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/model"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(StoresvcService) StoresvcService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type StoresvcService interface {
	// [method=post,expose=true,router=stores]
	Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Store, err error)
	// [method=post,expose=true,router=sync]
	Sync(ctx context.Context) (err error)
}

// the concrete implementation of service interface
type stubStoresvcService struct {
	logger log.Logger
	repo   model.StoreRepository
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.StoreRepository, logger log.Logger) (s StoresvcService) {
	var svc StoresvcService
	{
		svc = &stubStoresvcService{
			repo:   repo,
			logger: logger,
		}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Query
func (st *stubStoresvcService) Query(ctx context.Context, centerLng, centerLat, neLng, neLat, _, _, swLng, swLat, _, _ float64, max uint64) (items []model.Store, err error) {
	items, err = st.repo.Query(ctx, centerLng, centerLat, swLng, neLng, swLat, neLat, max)
	if err != nil {
		return items, err
	}
	return
}

// Implement the business logic of Sync
func (st *stubStoresvcService) Sync(ctx context.Context) (err error) {
	//resp, err := http.Get("https://raw.githubusercontent.com/kiang/pharmacies/master/json/points.json")
	//if err != nil {
	//	return errors.Wrap(errors.New("fetup point fail"), err)
	//}
	//defer resp.Body.Close()
	//
	//var req Collection
	//err = json.NewDecoder(resp.Body).Decode(&req)
	//if err != nil {
	//	msg := fmt.Sprintf("json.NewDecoder decode: %v", err)
	//	http.Error(w, msg, http.StatusBadRequest)
	//	return
	//}

	return err
}
