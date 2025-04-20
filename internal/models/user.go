package models

import (
	"time"
)

type User struct {
	UserId    int    `json:"user_id" gorm:"primarykey"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	RoleId    int    `json:"role_id"`

	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"-" gorm:"column:password_hash"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
