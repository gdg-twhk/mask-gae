package model

import (
	"context"
	"encoding/json"
	"time"
)

type Store struct {
	Id        string    `json:"id" db:"id"`
	Distance  float64   `json:"distance" db:"distance"`
	Name      string    `json:"name" db:"name"`
	Phone     string    `json:"phone" db:"phone"`
	Address   string    `json:"address" db:"address"`
	MaskAdult uint64    `json:"maskAdult" db:"mask_adult"`
	MaskChild uint64    `json:"maskChild" db:"mask_child"`
	Updated   time.Time `json:"updated" db:"updated"`
	Available string    `json:"available" db:"available"`
	Note      string    `json:"note" db:"note"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Latitude  float64   `json:"latitude" db:"latitude"`
}

func (p *Store) MarshalJSON() ([]byte, error) {
	type Alias Store

	return json.Marshal(&struct {
		*Alias
		Updated string `json:"updated"`
	}{
		Alias:   (*Alias)(p),
		Updated: p.Updated.Format("2006-01-02T15:04:05+08:00"),
	})
}

type StoreRepository interface {
	Query(context.Context, float64, float64, float64, float64, float64, float64, uint64) ([]Store, error)
	Insert(context.Context, Store) error
}
