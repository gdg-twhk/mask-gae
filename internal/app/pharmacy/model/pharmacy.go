package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lib/pq"

	"github.com/cage1016/mask/internal/pkg/util"
)

var location *time.Location

type Pharmacies []Pharmacy

func (p Pharmacies) Split(limit int) [][]Pharmacy {
	var chunk []Pharmacy
	chunks := make([][]Pharmacy, 0, len(p)/limit+1)
	for len(p) >= limit {
		chunk, p = p[:limit], p[limit:]
		chunks = append(chunks, chunk)
	}
	if len(p) > 0 {
		chunks = append(chunks, p[:len(p)])
	}
	return chunks
}

type Pharmacy struct {
	Id             string       `json:"id" db:"id"`
	Distance       float64      `json:"distance" db:"distance"`
	Name           string       `json:"name" db:"name"`
	Phone          string       `json:"phone" db:"phone"`
	Address        string       `json:"address" db:"address"`
	MaskAdult      uint64       `json:"maskAdult" db:"mask_adult"`
	MaskChild      uint64       `json:"maskChild" db:"mask_child"`
	Updated        *pq.NullTime `json:"updated" db:"updated"`
	Available      string       `json:"available" db:"available"`
	CustomNote     string       `json:"customNote" db:"custom_note"`
	Website        string       `json:"website" db:"website"`
	Note           string       `json:"note" db:"note"`
	Longitude      float64      `json:"longitude" db:"longitude"`
	Latitude       float64      `json:"latitude" db:"latitude"`
	ServicePeriods string       `json:"servicePeriods" db:"service_periods"`
	ServiceNote    string       `json:"serviceNote" db:"service_note"`
	County         string       `json:"county" db:"county"`
	Town           string       `json:"town" db:"town"`
	Cunli          string       `json:"cunli" db:"cunli"`
}

func (p *Pharmacy) MarshalJSON() ([]byte, error) {
	type Alias Pharmacy

	return json.Marshal(&struct {
		*Alias
		Updated string `json:"updated"`
	}{
		Alias: (*Alias)(p),
		Updated: func() string {
			if p.Updated != nil {
				return p.Updated.Time.In(util.Location).Format(time.RFC3339)
			}
			return ""
		}(),
	})
}

type PharmacyRepository interface {
	Query(context.Context, float64, float64, float64, float64, float64, float64, uint64) ([]Pharmacy, error)
	Insert(context.Context, string, [][]Pharmacy) error
	FootGun(ctx context.Context) error
}
