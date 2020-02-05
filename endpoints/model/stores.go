package model

type Stores struct {
	Name        string
	Phone       string
	Address     string
	MaskAdult   uint64
	MaskChild   uint64
	Updated     string
	Available   string
	Note        string
	Coordinates []float64
}
