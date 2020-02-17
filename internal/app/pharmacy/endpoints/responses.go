package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/app/pharmacy/service"
	"github.com/cage1016/mask/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*QueryResponse)(nil)

	_ httptransport.StatusCoder = (*QueryResponse)(nil)

	_ httptransport.Headerer = (*SyncResponse)(nil)

	_ httptransport.StatusCoder = (*SyncResponse)(nil)

	_ httptransport.Headerer = (*SyncHandlerResponse)(nil)

	_ httptransport.StatusCoder = (*SyncHandlerResponse)(nil)
)

// QueryResponse collects the response values for the Query method.
type QueryResponse struct {
	Items []model.Pharmacy `json:"items"`
	Err   error            `json:"err,omitempty"`
}

func (r QueryResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r QueryResponse) Headers() http.Header {
	return http.Header{}
}

func (r QueryResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}

func (r QueryResponse) ResponseOld() interface{} {
	return r.Items
}

// SyncResponse collects the response values for the Sync method.
type SyncResponse struct {
	Err error `json:"err,omitempty"`
}

func (r SyncResponse) StatusCode() int {
	return http.StatusOK
}

func (r SyncResponse) Headers() http.Header {
	return http.Header{}
}

func (r SyncResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// SyncHandlerResponse collects the response values for the SyncHandler method.
type SyncHandlerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r SyncHandlerResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r SyncHandlerResponse) Headers() http.Header {
	return http.Header{}
}

func (r SyncHandlerResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// FootGunResponse collects the response values for the FootGun method.
type FootGunResponse struct {
	Err error `json:"err,omitempty"`
}

func (r FootGunResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r FootGunResponse) Headers() http.Header {
	return http.Header{}
}

func (r FootGunResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}
