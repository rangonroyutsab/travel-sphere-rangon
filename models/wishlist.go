package models

import "time"

const (
	StatusPlanned = "Planned"
	StatusVisited = "Visited"
)

type WishlistItem struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	CountryName string    `json:"country_name"`
	Note        string    `json:"note"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
