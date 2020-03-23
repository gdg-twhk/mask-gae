package responses

type DataRes struct {
	APIVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type Responser interface {
	Response() interface{}
}

type Paging struct {
	CurrentItemCount int64 `json:"currentItemCount"`
	ItemsPage        int64 `json:"itemsPage"`
	StartIndex       int64 `json:"startIndex"`
	TotalItems       int64 `json:"totalItems"`
}