package models

import "time"

type Scope struct {
	ScopeId     int    `json:"scope_id" gorm:"primarykey"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Active    bool      `json:"active"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
