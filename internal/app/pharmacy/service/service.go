package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/storage"
	"github.com/go-kit/kit/log"
	"github.com/lib/pq"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/level"
	"github.com/cage1016/mask/internal/pkg/util"
)

var l log.Logger

var (
	ErrInvalidTask        = errors.New("Bad Request - Invalid Task")
	ErrCloudTaskNewClient = errors.New("task new client failed")
	ErrTaskCreatFailed    = errors.New("task create failed")
	ErrMalformedEntity    = errors.New("malformed entity specification")
	ErrStorageNewClient   = errors.New("storage new client failed")
	ErrStorageReadObject  = errors.New("storage read object failed")
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(PharmacyService) PharmacyService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
type PharmacyService interface {
	// [method=post,expose=true,router=api/pharmacies]
	Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Pharmacy, err error)
	// [method=post,expose=true,router=api/pharmacies/footgun]
	FootGun(ctx context.Context) (err error)
	// [method=post,expose=true,router=api/pharmacies/sync]
	Sync(ctx context.Context) (err error)
	// [method=post,expose=true,router=api/pharmacies/sync_handler]
	SyncHandler(ctx context.Context, queueName string, taskName string) (err error)
}

// the concrete implementation of service interface
type stubPharmacyService struct {
	logger           log.Logger
	repo             model.PharmacyRepository
	projectID        string
	locationID       string
	queueID          string
	bucketID         string
	pointsObjectName string
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.PharmacyRepository, projectID, LocationID, QueueID, BucketID, PointsObjectName string, logger log.Logger) (s PharmacyService) {
	var svc PharmacyService
	{
		l = logger
		svc = &stubPharmacyService{repo: repo, logger: logger, projectID: projectID, locationID: LocationID, queueID: QueueID, bucketID: BucketID, pointsObjectName: PointsObjectName}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Query
func (st *stubPharmacyService) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, _ float64, _ float64, swLng float64, swLat float64, _ float64, _ float64, max uint64) (items []model.Pharmacy, err error) {
	return st.repo.Query(ctx, centerLng, centerLat, swLng, neLng, swLat, neLat, max)
}

// Implement the business logic of FootGun
func (ph *stubPharmacyService) FootGun(ctx context.Context) (err error) {
	return ph.repo.FootGun(ctx)
}

// Implement the business logic of Sync
func (st *stubPharmacyService) Sync(ctx context.Context) (err error) {
	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		level.Error(st.logger).Log("method", "cloudtasks.NewClient", "err", err)
		return errors.Wrap(ErrCloudTaskNewClient, err)
	}
	defer taskClient.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", st.projectID, st.locationID, st.queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#AppEngineHttpRequest
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/api/pharmacies/sync_handler",
				},
			},
		},
	}

	createdTask, err := taskClient.CreateTask(ctx, req)
	if err != nil {
		level.Error(st.logger).Log("method", "taskClient.CreateTask", "err", err)
		return errors.Wrap(ErrTaskCreatFailed, err)
	}

	level.Info(st.logger).Log("method", "taskClient.CreateTask", "task", createdTask)
	return nil
}

// Implement the business logic of SyncHandler
func (st *stubPharmacyService) SyncHandler(ctx context.Context, queueName string, taskName string) (err error) {
	level.Info(st.logger).Log("method", "SyncHandler start", "queue", queueName, "task", taskName)
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		level.Error(st.logger).Log("method", "storage.NewClient", "err", err)
		return errors.Wrap(ErrStorageNewClient, err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(st.bucketID).Object(st.pointsObjectName).NewReader(ctx)
	if err != nil {
		level.Error(st.logger).Log("method", `client.Bucket(BucketID).Object("/points.json").NewReader(ctx)`, "err", err)
		return errors.Wrap(ErrStorageReadObject, err)
	}
	defer rc.Close()

	var req Collection
	err = json.NewDecoder(rc).Decode(&req)
	if err != nil {
		level.Error(st.logger).Log("method", "json.NewDecoder(rc).Decode(&req)", "err", err)
		return errors.Wrap(ErrMalformedEntity, err)
	}

	var updated string
	pharmacies := make(model.Pharmacies, len(req.Features))
	for i, f := range req.Features {
		if updated == "" {
			if f.Properties.Updated.Valid {
				updated = f.Properties.Updated.Time.In(util.Location).Format("2006_0102_150405")
			}
		}
		pharmacy := model.Pharmacy{
			Id:             f.Properties.Id,
			Name:           f.Properties.Name,
			Phone:          f.Properties.Phone,
			Address:        f.Properties.Address,
			MaskAdult:      f.Properties.MaskAdult,
			MaskChild:      f.Properties.MaskChild,
			Updated:        f.Properties.Updated,
			Available:      f.Properties.Available,
			CustomNote:     f.Properties.CustomNote,
			Website:        f.Properties.Website,
			Note:           f.Properties.Note,
			Longitude:      f.Geometry.Coordinates[0],
			Latitude:       f.Geometry.Coordinates[1],
			ServicePeriods: f.Properties.ServicePeriods,
			ServiceNote:    f.Properties.ServiceNote,
			County:         f.Properties.County,
			Town:           f.Properties.Town,
			Cunli:          f.Properties.Cunli,
		}
		pharmacies[i] = pharmacy
	}

	err = st.repo.Insert(ctx, updated, pharmacies.Split(500))
	if err != nil {
		level.Error(st.logger).Log("method", "st.repo.Insert", "queue", queueName, "task", taskName, "err", err)
	} else {
		level.Info(st.logger).Log("method", "st.repo.Insert", "queue", queueName, "task", taskName)
	}
	return err
}

type Properties struct {
	Id             string       `json:"id"`
	Name           string       `json:"name"`
	Phone          string       `json:"phone"`
	Address        string       `json:"address"`
	MaskAdult      uint64       `json:"mask_adult"`
	MaskChild      uint64       `json:"mask_child"`
	Updated        *pq.NullTime `json:"updated"`
	Note           string       `json:"note"`
	Available      string       `json:"available"`
	CustomNote     string       `json:"custom_note"`
	Website        string       `json:"website"`
	ServicePeriods string       `json:"service_periods"`
	ServiceNote    string       `json:"service_note"`
	County         string       `json:"county"`
	Town           string       `json:"town"`
	Cunli          string       `json:"cunli"`
}

// UnmarshalJSON means to unmarshal json to object.
func (p *Properties) UnmarshalJSON(data []byte) error {
	type Alias Properties

	pr := &struct {
		Updated string `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &pr); err != nil {
		level.Warn(l).Log("method", "json.Unmarshal", "err", err, "properties", fmt.Sprintf("%+v", p))
		return nil
	}

	if pr.Updated != "" {
		expired, err := time.ParseInLocation("2006/01/02 15:04:05", pr.Updated, util.Location)
		if err != nil {
			return err
		}

		p.Updated = &pq.NullTime{
			Time:  expired,
			Valid: true,
		}
	}

	return nil
}

type Features struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

type Collection struct {
	Type     string     `json:"type"`
	Features []Features `json:"features"`
}
