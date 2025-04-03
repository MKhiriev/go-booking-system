package models

import (
	"encoding/json"
	"log"
	"time"
)

// Room TODO add getters/setters
type Room struct {
	RoomId   int    `json:"room_id"`
	Number   string `json:"number"`
	Capacity int    `json:"capacity"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func ParseRoom(data []byte) (*Room, error) {
	room := &Room{}
	err := json.Unmarshal(data, room)
	if err != nil {
		log.Println("Models.ParseRoom error: ", err)
		return room, err
	}
	return room, nil
}
