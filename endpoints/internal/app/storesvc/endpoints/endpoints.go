package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/model"
	"github.com/cage1016/mask/internal/app/storesvc/service"
)

// Endpoints collects all of the endpoints that compose the storesvc service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	QueryEndpoint endpoint.Endpoint
	SyncEndpoint  endpoint.Endpoint
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.StoresvcService, logger log.Logger) (ep Endpoints) {
	var queryEndpoint endpoint.Endpoint
	{
		method := "query"
		queryEndpoint = MakeQueryEndpoint(svc)
		queryEndpoint = LoggingMiddleware(log.With(logger, "method", method))(queryEndpoint)
		ep.QueryEndpoint = queryEndpoint
	}

	var syncEndpoint endpoint.Endpoint
	{
		method := "sync"
		syncEndpoint = MakeSyncEndpoint(svc)
		syncEndpoint = LoggingMiddleware(log.With(logger, "method", method))(syncEndpoint)
		ep.SyncEndpoint = syncEndpoint
	}

	return ep
}

// MakeQueryEndpoint returns an endpoint that invokes Query on the service.
// Primarily useful in a server.
func MakeQueryEndpoint(svc service.StoresvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(QueryRequest)
		if err := req.validate(); err != nil {
			return QueryResponse{}, err
		}
		stores, err := svc.Query(ctx, req.Center.Lng, req.Center.Lat, req.Bounds.Ne.Lng, req.Bounds.Ne.Lat, req.Bounds.Se.Lng, req.Bounds.Se.Lat, req.Bounds.Sw.Lng, req.Bounds.Sw.Lat, req.Bounds.Nw.Lng, req.Bounds.Nw.Lat, req.Max)
		return QueryResponse{Items: stores}, err
	}
}

// Query implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Store, err error) {
	resp, err := e.QueryEndpoint(ctx, QueryRequest{Center: LatLng{
		Lat: centerLat,
		Lng: centerLng,
	}, Bounds: Bounds{
		Ne: LatLng{neLat, neLng},
		Se: LatLng{seLat, seLng},
		Sw: LatLng{swLat, swLng},
		Nw: LatLng{nwLat, nwLng},
	}, Max: max})
	if err != nil {
		return
	}
	response := resp.(QueryResponse)
	return response.Items, nil
}

// MakeSyncEndpoint returns an endpoint that invokes Sync on the service.
// Primarily useful in a server.
func MakeSyncEndpoint(svc service.StoresvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.Sync(ctx)
		return SyncResponse{}, err
	}
}

// Sync implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Sync(ctx context.Context) (err error) {
	resp, err := e.SyncEndpoint(ctx, SyncRequest{})
	if err != nil {
		return
	}
	_ = resp.(SyncResponse)
	return nil
}
