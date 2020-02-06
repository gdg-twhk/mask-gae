package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/mask/internal/app/model"
	"github.com/cage1016/mask/internal/app/storesvc/service"
	"github.com/cage1016/mask/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*QueryResponse)(nil)

	_ httptransport.StatusCoder = (*QueryResponse)(nil)

	_ httptransport.Headerer = (*SyncResponse)(nil)

	_ httptransport.StatusCoder = (*SyncResponse)(nil)
)

// QueryResponse collects the response values for the Query method.
type QueryResponse struct {
	Items []model.Store `json:"items"`
	Err   error         `json:"err,omitempty"`
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

// SyncResponse collects the response values for the Sync method.
type SyncResponse struct {
	Err error `json:"err"`
}

func (r SyncResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r SyncResponse) Headers() http.Header {
	return http.Header{}
}

func (r SyncResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}
