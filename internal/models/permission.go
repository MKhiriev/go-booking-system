package models

import "time"

type Permission struct {
	RoleId  int `json:"role_id"`
	RouteId int `json:"route_id"`
	ScopeId int `json:"scope_id"`

	Active bool `json:"active"`
	// TODO - uncomment after creating Repositories, Services, Handlers for Scope, Permission, Route
	// CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
