package models

type Attraction struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Kind     string  `json:"kind"`
	Distance float64 `json:"distance"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}
