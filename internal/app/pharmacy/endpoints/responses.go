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

	_ httptransport.Headerer = (*TickerUpdateResponse)(nil)

	_ httptransport.StatusCoder = (*TickerUpdateResponse)(nil)
)

// QueryResponse collects the response values for the Query method.
type QueryResponse struct {
	Items []model.Pharmacy `json:"items"`
	Err   error            `json:"-"`
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

// TickerUpdateResponse collects the response values for the TickerUpdate method.
type TickerUpdateResponse struct {
	Err error `json:"err"`
}

func (r TickerUpdateResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r TickerUpdateResponse) Headers() http.Header {
	return http.Header{}
}

func (r TickerUpdateResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}
