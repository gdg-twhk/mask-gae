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

// SyncRequest collects the request parameters for the Sync method.
type SyncRequest struct {
}

func (r SyncRequest) validate() error {
	return nil // TBA
}

// SyncHandlerRequest collects the request parameters for the SyncHandler method.
type SyncHandlerRequest struct {
	QueueName string `json:"queue_name"`
	TaskName  string `json:"task_name"`
}

func (r SyncHandlerRequest) validate() error {
	return nil // TBA
}
