package model

import "time"

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
