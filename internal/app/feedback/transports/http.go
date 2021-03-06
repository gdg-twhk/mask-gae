package transports

import (
	"context"
	"encoding/json"
	"github.com/rs/cors"
	"strconv"
	"time"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"io"
	"net/http"

	"github.com/cage1016/mask/internal/app/feedback/endpoints"
	"github.com/cage1016/mask/internal/app/feedback/service"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/responses"
)

const (
	contentType string = "application/json"

	defOffset = 0
	defLimit  = 10
)

// ShowFeedback godoc
// @Summary feedback options
// @Description The endpoint for Retailbase to fetch feedback options
// @Tags feedback
// @Accept json
// @Produce json
// @Success 200 {object} endpoints.OptionsResponse
// @Failure 400 {object} responses.ErrorRes
// @Failure 500 {object} responses.ErrorRes
// @Router /api/feedback/options [get]
func OptionsHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, logger log.Logger) {
	m.Get("/api/feedback/options", httptransport.NewServer(
		endpoints.OptionsEndpoint,
		decodeHTTPOptionsRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))

}

// ShowFeedback godoc
// @Summary specific pharmacy feedbacks
// @Description The endpoint for Retailbase to fetch specific pharmacy feedbacks
// @Tags feedback
// @Accept json
// @Produce json
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   date      query    string     true       "date, yyyy_mmdd"
// @Param pharmacy_id path string true "Pharmacy ID"
// @Success 200 {object} model.FeedbackItemPage
// @Failure 400 {object} responses.ErrorRes
// @Failure 500 {object} responses.ErrorRes
// @Router /api/feedback/pharmacies/{pharmacy_id} [get]
func PharmacyFeedBacksHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, logger log.Logger) {
	m.Get("/api/feedback/pharmacies/:pharmacy_id", httptransport.NewServer(
		endpoints.PharmacyFeedBacksEndpoint,
		decodeHTTPPharmacyFeedBacksRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))

}

// ShowFeedback godoc
// @Summary specific user feedbacks
// @Description The endpoint for Retailbase to fetch specific user feedbacks
// @Tags feedback
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param   offset     query    int     true        "Offset"
// @Param   limit      query    int     true        "limit"
// @Param   date      query    string     true       "date, yyyy_mmdd"
// @Success 200 {object} model.FeedbackItemPage
// @Failure 400 {object} responses.ErrorRes
// @Failure 500 {object} responses.ErrorRes
// @Router /api/feedback/users/{user_id} [get]
func UserFeedBacksHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, logger log.Logger) {
	m.Get("/api/feedback/users/:user_id", httptransport.NewServer(
		endpoints.UserFeedBacksEndpoint,
		decodeHTTPUserFeedBacksRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))

}

// ShowFeedback godoc
// @Summary specific pharmacy feedbacks
// @Description The endpoint for Retailbase to fetch specific pharmacy feedbacks
// @Tags feedback
// @Accept json
// @Produce json
// @Param user body endpoints.FeedBackRequest true "Feedback"
// @Success 200 {object} endpoints.FeedBackResponse
// @Failure 400 {object} responses.ErrorRes
// @Failure 500 {object} responses.ErrorRes
// @Router /api/feedback [post]
func FeedBackHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, logger log.Logger) {
	m.Post("/api/feedback", httptransport.NewServer(
		endpoints.FeedBackEndpoint,
		decodeHTTPFeedBackRequest,
		encodeJSONResponse,
		append(options, httptransport.ServerBefore(kitjwt.HTTPToContext()))...,
	))
}

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints, logger log.Logger) http.Handler { // Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(httpEncodeError),
		httptransport.ServerErrorLogger(logger),
	}

	m := bone.New()
	OptionsHandler(m, endpoints, options, logger)
	PharmacyFeedBacksHandler(m, endpoints, options, logger)
	UserFeedBacksHandler(m, endpoints, options, logger)
	FeedBackHandler(m, endpoints, options, logger)
	return cors.AllowAll().Handler(m)
}

// decodeHTTPOptionsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPOptionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.OptionsRequest
	return req, nil
}

// decodeHTTPPharmacyFeedBacksRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPPharmacyFeedBacksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.PharmacyFeedBacksRequest
	req.PharmacyID = bone.GetValue(r, "pharmacy_id")
	s := bone.GetQuery(r, "date")
	if len(s) > 0 {
		req.Date = s[0]
	} else {
		req.Date = time.Now().Format(service.QueryDatefmt)
	}

	var err error
	req.Offset, err = readUintQuery(r, "offset", defOffset)
	if err != nil {
		return nil, err
	}

	req.Limit, err = readUintQuery(r, "limit", defLimit)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// decodeHTTPUserFeedBacksRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPUserFeedBacksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UserFeedBacksRequest
	req.UserID = bone.GetValue(r, "user_id")
	s := bone.GetQuery(r, "date")
	if len(s) > 0 {
		req.Date = s[0]
	} else {
		req.Date = time.Now().Format(service.QueryDatefmt)
	}

	var err error
	req.Offset, err = readUintQuery(r, "offset", 0)
	if err != nil {
		return nil, err
	}

	req.Limit, err = readUintQuery(r, "limit", 10)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// decodeHTTPFeedBackRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPFeedBackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.FeedBackRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func httpEncodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	var message string
	var errs []errors.Errors
	w.Header().Set("Content-Type", contentType)

	// HTTP
	switch errorVal := err.(type) {
	case errors.Error:
		switch {
		case errors.Contains(errorVal, service.ErrMalformedEntity):
			code = http.StatusBadRequest
		}

		if errorVal.Msg() != "" {
			message, errs = errorVal.Msg(), errorVal.Errors()
		}
	default:
		switch err {
		case io.ErrUnexpectedEOF, io.EOF:
			code = http.StatusBadRequest
		case kitjwt.ErrTokenContextMissing:
			code = http.StatusUnauthorized
		default:
			switch err.(type) {
			case *json.SyntaxError, *json.UnmarshalTypeError:
				code = http.StatusBadRequest
			}
		}

		errs = errors.FromError(err.Error())
		message = errs[0].Message
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(responses.ErrorRes{responses.ErrorResItem{code, message, errs}})
}

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	if ar, ok := response.(responses.Responser); ok {
		return json.NewEncoder(w).Encode(ar.Response())
	}

	return json.NewEncoder(w).Encode(response)
}

func readUintQuery(r *http.Request, key string, def uint64) (uint64, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return 0, service.ErrInvalidQueryParams
	}

	if len(vals) == 0 {
		return def, nil
	}

	strval := vals[0]
	val, err := strconv.ParseUint(strval, 10, 64)
	if err != nil {
		return 0, service.ErrInvalidQueryParams
	}

	return val, nil
}
