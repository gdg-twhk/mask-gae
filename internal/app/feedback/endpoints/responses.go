package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/mask/internal/app/feedback/model"
	"github.com/cage1016/mask/internal/app/feedback/service"
	"github.com/cage1016/mask/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*OptionsResponse)(nil)

	_ httptransport.StatusCoder = (*OptionsResponse)(nil)

	_ httptransport.Headerer = (*PharmacyFeedBacksResponse)(nil)

	_ httptransport.StatusCoder = (*PharmacyFeedBacksResponse)(nil)

	_ httptransport.Headerer = (*UserFeedBacksResponse)(nil)

	_ httptransport.StatusCoder = (*UserFeedBacksResponse)(nil)

	_ httptransport.Headerer = (*FeedBackResponse)(nil)

	_ httptransport.StatusCoder = (*FeedBackResponse)(nil)
)

// OptionsResponse collects the response values for the Options method.
type OptionsResponse struct {
	Items []model.Option `json:"items"`
	Err   error          `json:"err,omitempty"`
}

func (r OptionsResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r OptionsResponse) Headers() http.Header {
	return http.Header{}
}

func (r OptionsResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}

// PharmacyFeedBacksResponse collects the response values for the PharmacyFeedBacks method.
type PharmacyFeedBacksResponse struct {
	Res model.FeedbackItemPage `json:"items"`
	Err error                  `json:"err,omitempty"`
}

func (r PharmacyFeedBacksResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r PharmacyFeedBacksResponse) Headers() http.Header {
	return http.Header{}
}

func (r PharmacyFeedBacksResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// UserFeedBacksResponse collects the response values for the UserFeedBacks method.
type UserFeedBacksResponse struct {
	Res model.FeedbackItemPage `json:"res"`
	Err error                  `json:"err,omitempty"`
}

func (r UserFeedBacksResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r UserFeedBacksResponse) Headers() http.Header {
	return http.Header{}
}

func (r UserFeedBacksResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// FeedBackResponse collects the response values for the InsertFeedBack method.
type FeedBackResponse struct {
	Err error  `json:"err,omitempty"`
	ID  string `json:"id"`
}

func (r FeedBackResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r FeedBackResponse) Headers() http.Header {
	return http.Header{}
}

func (r FeedBackResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}
