package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/cage1016/mask/model"
)

var (
	ProjectID  = os.Getenv("PROJECT_ID")
	LocationID = "asia-east2"
	QueueID    = "sync-points-queue2"
	BucketID   = "mask-9999-pharmacies"
)

var db *sqlx.DB

func init() {
	connectionName := os.Getenv("CLOUDSQL_CONNECTION_NAME")
	var err error
	// PostgreSQL Connection, uncomment to use.
	// connection string format: user=USER password=PASSWORD host=/cloudsql/PROJECT_ID:REGION_ID:INSTANCE_ID/[ dbname=DB_NAME]
	// dbURI := fmt.Sprintf("user=%s password=%s host=/cloudsql/%s dbname=%s", user, password, connectionName, dbName)
	// conn, err := sql.Open("postgres", dbURI)
	dbURI := fmt.Sprintf("user=%s password=%s host=/cloudsql/%s dbname=%s sslmode=disable", "postgres", "password", connectionName, "stores")
	//dbURI = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", "postgres", "password", "localhost", "stores")
	db, err = sqlx.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Request struct {
	Center LatLng `json:"center"`
	Bounds struct {
		Ne LatLng `json:"ne"`
		Se LatLng `json:"se"`
		Sw LatLng `json:"sw"`
		Nw LatLng `json:"nw"`
	} `json:"bounds"`
	Max uint64 `json:"max"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to mask-endpoints")
}

func StoresHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		code := http.StatusBadRequest
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := http.StatusOK
	w.WriteHeader(code)

	q := `SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from stores where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`

	stroes := []model.Store{}
	if err := db.Select(&stroes, q, req.Center.Lng, req.Center.Lat, req.Bounds.Sw.Lng, req.Bounds.Ne.Lng, req.Bounds.Sw.Lat, req.Bounds.Ne.Lat, req.Max); err != nil {
		msg := fmt.Sprintf("db.Query: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	val, _ := json.Marshal(stroes)
	w.Write(val)
}

type Properties struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	MaskAdult uint64    `json:"mask_adult"`
	MaskChild uint64    `json:"mask_child"`
	Updated   time.Time `json:"updated"`
	Note      string    `json:"note"`
	Available string    `json:"available"`
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
		return nil
	}

	if pr.Updated != "" {
		expired, err := time.Parse("2006/01/02 15:04:05", pr.Updated)
		if err != nil {
			return err
		}

		p.Updated = expired
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

func SyncQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
		return
	}
	defer taskClient.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", ProjectID, LocationID, QueueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#AppEngineHttpRequest
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/api/sync_handler",
				},
			},
		},
	}

	task, err := taskClient.CreateTask(ctx, req)
	if err != nil {
		log.Fatalf("taskClient.CreateTask fail: %s", err.Error())
		return
	}
	log.Println(ctx, "order email task added: %+v", task)

	var encjson string
	var status int
	if err != nil {
		encjson = ConvertToJson(err)
		status = http.StatusBadRequest
	} else {
		encjson = ConvertToJson(task)
		status = http.StatusOK
	}

	if status != http.StatusOK {
		http.Error(w, string(encjson), status)
	} else {
		fmt.Fprintln(w, encjson)
	}
}

func ConvertToJson(v interface{}) string {
	encjson, _ := json.Marshal(v)
	return string(encjson)
}

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	t, ok := r.Header["X-Appengine-Taskname"]
	if !ok || len(t[0]) == 0 {
		// You may use the presence of the X-Appengine-Taskname header to validate
		// the request comes from Cloud Tasks.
		log.Println("Invalid Task: No X-Appengine-Taskname request header found")
		http.Error(w, "Bad Request - Invalid Task", http.StatusBadRequest)
		return
	}
	taskName := t[0]

	// Pull useful headers from Task request.
	q, ok := r.Header["X-Appengine-Queuename"]
	queueName := ""
	if ok {
		queueName = q[0]
	}

	log.Printf("SyncHandler start task queue(%s), task name(%s)", taskName, queueName)
	// Creates a client.
	ctx := r.Context()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	log.Println("storage.NewClient ok")
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(BucketID).Object("/points.json").NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
	defer rc.Close()
	log.Println("read points.json ok")

	var req Collection
	err = json.NewDecoder(rc).Decode(&req)
	if err != nil {
		msg := fmt.Sprintf("json.NewDecoder decode: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	log.Println("parse points.json ok")
	log.Println(len(req.Features))


	if _, err :=db.Exec(`delete from stores`); err != nil {
		log.Fatalf("delete from stores err: %v", err)
	}

	ql := `INSERT INTO public.stores (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude, updated)
		VALUES (:id, :name, :phone, :address, :mask_adult, :mask_child, :available, :note, :longitude, :latitude, :updated)`

	log.Println("features counts ", len(req.Features))

	for _, f := range req.Features {
		store := model.Store{
			Id:        f.Properties.Id,
			Name:      f.Properties.Name,
			Phone:     f.Properties.Phone,
			Address:   f.Properties.Address,
			MaskAdult: f.Properties.MaskAdult,
			MaskChild: f.Properties.MaskChild,
			Updated:   f.Properties.Updated,
			Available: f.Properties.Available,
			Note:      f.Properties.Note,
			Longitude: f.Geometry.Coordinates[0],
			Latitude:  f.Geometry.Coordinates[1],
		}

		if _, err := db.NamedExec(ql, store); err != nil {
			log.Fatalf("db.NamedExec: %v", err)
		}
	}
	log.Printf("sync handler done task queue(%s), task name(%s)", taskName, queueName)
	fmt.Fprintf(w, "sync handler done task queue(%s), task name(%s)", taskName, queueName)
}
