package models

import "time"

type Route struct {
	RouteId     int    `json:"route_id" gorm:"primarykey"`
	URL         string `json:"url"`
	Description string `json:"description"`

	Active    bool      `json:"active"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
