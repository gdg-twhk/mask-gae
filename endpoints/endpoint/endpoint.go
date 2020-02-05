package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/cage1016/mask/model"
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
	//dbURI := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", "postgres", "password", "localhost", "stores")
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
	if err := db.Select(&stroes, q, req.Center.Lng, req.Center.Lat, req.Bounds.Sw.Lng, req.Bounds.Ne.Lng, req.Bounds.Se.Lat, req.Bounds.Ne.Lat, req.Max); err != nil {
		msg := fmt.Sprintf("db.Query: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	val, _ := json.Marshal(stroes)
	w.Write(val)
}

type Properties struct {
	Id        uint64    `json:"id"`
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
		expired, err := time.Parse("2006/01/02 15:04", pr.Updated)
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

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://raw.githubusercontent.com/kiang/pharmacies/master/json/points.json")
	defer resp.Body.Close()

	var req Collection
	_ = json.NewDecoder(resp.Body).Decode(&req)

	q := `INSERT INTO public.stores (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude,
                           updated)
		VALUES (:id, :name, :phone, :address, :mask_adult, :mask_child, :available, :note, :longitude, :latitude, :updated)
		ON CONFLICT (id)
			DO UPDATE
			SET name       = :name,
				phone      = :phone,
				address    = :address,
				mask_adult = :mask_adult,
				mask_child = :mask_child,
				available  = :available,
				note       = :note,
				longitude  = :longitude,
				latitude   = :latitude,
				updated    = :updated;`

	for _, f := range req.Features {
		store := model.Store{
			Id:        fmt.Sprintf("%d", f.Properties.Id),
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

		if _, err := db.NamedExec(q, store); err != nil {
			log.Fatalf("db.NamedExec: %v", err)
		}
	}

	w.Write([]byte("sync done"))
}
