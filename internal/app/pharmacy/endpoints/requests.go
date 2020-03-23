package endpoints

type Request interface {
	validate() error
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Bounds struct {
	Ne LatLng `json:"ne"`
	Se LatLng `json:"se"`
	Sw LatLng `json:"sw"`
	Nw LatLng `json:"nw"`
}

// QueryRequest collects the request parameters for the Query method.
type QueryRequest struct {
	Center LatLng `json:"center"`
	Bounds Bounds `json:"bounds"`
	Max    uint64 `json:"max"`
}

func (r QueryRequest) validate() error {
	return nil // TBA
}

// TickerUpdateRequest collects the request parameters for the TickerUpdate method.
type TickerUpdateRequest struct {
}

func (r TickerUpdateRequest) validate() error {
	return nil // TBA
}
