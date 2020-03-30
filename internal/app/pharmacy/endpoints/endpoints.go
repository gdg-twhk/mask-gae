package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/app/pharmacy/service"
)

// Endpoints collects all of the endpoints that compose the pharmacy service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	QueryEndpoint        endpoint.Endpoint `json:""`
	TickerUpdateEndpoint endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.PharmacyService, logger log.Logger) (ep Endpoints) {
	var queryEndpoint endpoint.Endpoint
	{
		method := "query"
		queryEndpoint = MakeQueryEndpoint(svc)
		queryEndpoint = LoggingMiddleware(log.With(logger, "method", method))(queryEndpoint)
		ep.QueryEndpoint = queryEndpoint
	}

	var tickerUpdateEndpoint endpoint.Endpoint
	{
		method := "tickerUpdate"
		tickerUpdateEndpoint = MakeTickerUpdateEndpoint(svc)
		tickerUpdateEndpoint = LoggingMiddleware(log.With(logger, "method", method))(tickerUpdateEndpoint)
		ep.TickerUpdateEndpoint = tickerUpdateEndpoint
	}

	return ep
}

// MakeQueryEndpoint returns an endpoint that invokes Query on the service.
// Primarily useful in a server.
func MakeQueryEndpoint(svc service.PharmacyService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(QueryRequest)
		if err := req.validate(); err != nil {
			return QueryResponse{}, err
		}
		pharmacies, err := svc.Query(ctx, req.Center.Lng, req.Center.Lat, req.Bounds.Ne.Lng, req.Bounds.Ne.Lat, req.Bounds.Se.Lng, req.Bounds.Se.Lat, req.Bounds.Sw.Lng, req.Bounds.Sw.Lat, req.Bounds.Nw.Lng, req.Bounds.Nw.Lat, req.Max)
		return QueryResponse{Items: pharmacies}, err
	}
}

// Query implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Pharmacy, err error) {
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

// MakeTickerUpdateEndpoint returns an endpoint that invokes TickerUpdate on the service.
// Primarily useful in a server.
func MakeTickerUpdateEndpoint(svc service.PharmacyService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.TickerUpdate(ctx)
		return TickerUpdateResponse{}, err
	}
}

// TickerUpdate implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) TickerUpdate(ctx context.Context) (err error) {
	resp, err := e.TickerUpdateEndpoint(ctx, TickerUpdateRequest{})
	if err != nil {
		return
	}
	_ = resp.(TickerUpdateResponse)
	return nil
}
