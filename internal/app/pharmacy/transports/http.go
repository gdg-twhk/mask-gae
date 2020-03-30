package transports

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/rs/cors"

	"github.com/cage1016/mask/internal/app/pharmacy/endpoints"
	"github.com/cage1016/mask/internal/app/pharmacy/service"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/responses"
)

const (
	contentType string = "application/json"
)

// ShowPharmacy godoc
// @Summary
// @Description The endpoint for Mask to fetch pharmacies
// @Tags pharmacy
// @Accept json
// @Produce json
// @Param query body endpoints.QueryRequest true "Fetch Pharmacies"
// @Success 200 {object} endpoints.QueryResponse
// @Failure 400 {object} responses.ErrorRes
// @Failure 500 {object} responses.ErrorRes
// @Router /api/pharmacies [post]
func QueryHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, logger log.Logger) {
	m.Post("/api/pharmacies", httptransport.NewServer(
		endpoints.QueryEndpoint,
		decodeHTTPQueryRequest,
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
	QueryHandler(m, endpoints, options, logger)
	m.GetFunc("/_ah/warmup", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("/_ah/warmup", "done")
	})
	return cors.AllowAll().Handler(m)
}

// decodeHTTPQueryRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPQueryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.QueryRequest
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
		case errors.Contains(errorVal, service.ErrMalformedEntity),
			errors.Contains(errorVal, service.ErrInvalidTask),
			errors.Contains(errorVal, service.ErrTaskCreatFailed):
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
