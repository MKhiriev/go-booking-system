package models

import (
	"time"
)

type Room struct {
	RoomId   int    `json:"room_id" gorm:"primarykey"`
	Number   string `json:"number"`
	Capacity int    `json:"capacity"`

	Active    bool      `json:"active"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
