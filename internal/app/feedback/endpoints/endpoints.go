package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/feedback/model"
	"github.com/cage1016/mask/internal/app/feedback/service"
)

// Endpoints collects all of the endpoints that compose the feedbacksvc service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	OptionsEndpoint           endpoint.Endpoint `json:""`
	PharmacyFeedBacksEndpoint endpoint.Endpoint `json:""`
	UserFeedBacksEndpoint     endpoint.Endpoint `json:""`
	FeedBackEndpoint          endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.FeedbacksvcService, logger log.Logger) (ep Endpoints) {
	var optionsEndpoint endpoint.Endpoint
	{
		method := "options"
		optionsEndpoint = MakeOptionsEndpoint(svc)
		optionsEndpoint = LoggingMiddleware(log.With(logger, "method", method))(optionsEndpoint)
		ep.OptionsEndpoint = optionsEndpoint
	}

	var storesEndpoint endpoint.Endpoint
	{
		method := "stores"
		storesEndpoint = MakePharmacyFeedBacksEndpoint(svc)
		storesEndpoint = LoggingMiddleware(log.With(logger, "method", method))(storesEndpoint)
		ep.PharmacyFeedBacksEndpoint = storesEndpoint
	}

	var usersEndpoint endpoint.Endpoint
	{
		method := "users"
		usersEndpoint = MakeUserFeedBacksEndpoint(svc)
		usersEndpoint = LoggingMiddleware(log.With(logger, "method", method))(usersEndpoint)
		ep.UserFeedBacksEndpoint = usersEndpoint
	}

	var feedBackEndpoint endpoint.Endpoint
	{
		method := "feedBack"
		feedBackEndpoint = MakeFeedBackEndpoint(svc)
		feedBackEndpoint = LoggingMiddleware(log.With(logger, "method", method))(feedBackEndpoint)
		ep.FeedBackEndpoint = feedBackEndpoint
	}

	return ep
}

// MakeOptionsEndpoint returns an endpoint that invokes Options on the service.
// Primarily useful in a server.
func MakeOptionsEndpoint(svc service.FeedbacksvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		items, err := svc.Options(ctx)
		return OptionsResponse{Items: items}, err
	}
}

// Options implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Options(ctx context.Context) (items []model.Option, err error) {
	resp, err := e.OptionsEndpoint(ctx, OptionsRequest{})
	if err != nil {
		return
	}
	response := resp.(OptionsResponse)
	return response.Items, nil
}

// MakePharmacyFeedBacksEndpoint returns an endpoint that invokes PharmacyFeedBacks on the service.
// Primarily useful in a server.
func MakePharmacyFeedBacksEndpoint(svc service.FeedbacksvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PharmacyFeedBacksRequest)
		if err := req.validate(); err != nil {
			return PharmacyFeedBacksResponse{}, err
		}
		res, err := svc.PharmacyFeedBacks(ctx, req.PharmacyID, req.Date, req.Offset, req.Limit)
		return PharmacyFeedBacksResponse{Res: res}, err
	}
}

// PharmacyFeedBacks implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) PharmacyFeedBacks(ctx context.Context, pharmacyID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	resp, err := e.PharmacyFeedBacksEndpoint(ctx, PharmacyFeedBacksRequest{PharmacyID: pharmacyID, Date: date, Offset: offset, Limit: limit})
	if err != nil {
		return
	}
	response := resp.(PharmacyFeedBacksResponse)
	return response.Res, nil
}

// MakeUserFeedBacksEndpoint returns an endpoint that invokes UserFeedBacks on the service.
// Primarily useful in a server.
func MakeUserFeedBacksEndpoint(svc service.FeedbacksvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserFeedBacksRequest)
		if err := req.validate(); err != nil {
			return UserFeedBacksResponse{}, err
		}
		res, err := svc.UserFeedBacks(ctx, req.UserID, req.Date, req.Offset, req.Limit)
		return UserFeedBacksResponse{Res: res}, err
	}
}

// UserFeedBacks implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) UserFeedBacks(ctx context.Context, userID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	resp, err := e.UserFeedBacksEndpoint(ctx, UserFeedBacksRequest{UserID: userID, Date: date, Offset: offset, Limit: limit})
	if err != nil {
		return
	}
	response := resp.(UserFeedBacksResponse)
	return response.Res, nil
}

// MakeFeedBackEndpoint returns an endpoint that invokes InsertFeedBack on the service.
// Primarily useful in a server.
func MakeFeedBackEndpoint(svc service.FeedbacksvcService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FeedBackRequest)
		if err := req.validate(); err != nil {
			return FeedBackResponse{}, err
		}
		id, err := svc.InsertFeedBack(ctx,
			req.UserID,
			req.PharmacyID,
			req.OptionID,
			req.Description,
			req.Longitude,
			req.Latitude,
		)
		return FeedBackResponse{ID: id}, err
	}
}

// InsertFeedBack implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) FeedBack(ctx context.Context, userID, pharmacyID, optionID, description string, Longitude, Latitude float64) (id string, err error) {
	resp, err := e.FeedBackEndpoint(ctx, FeedBackRequest{
		UserID:      userID,
		PharmacyID:  pharmacyID,
		OptionID:    optionID,
		Description: description,
		Longitude:   Longitude,
		Latitude:    Latitude,
	})
	if err != nil {
		return
	}
	response := resp.(FeedBackResponse)
	return response.ID, nil
}
