package models

import (
	"encoding/json"
	"log"
	"time"
)

type Booking struct {
	BookingId     int       `json:"booking_id"`
	UserId        int       `json:"user_id"`
	RoomId        int       `json:"room_id"`
	DateTimeStart time.Time `json:"datetime_start"`
	DateTimeEnd   time.Time `json:"datetime_end"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func ParseBooking(data []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(data, user)
	if err != nil {
		log.Println("Models.ParseBooking error: ", err)
		return user, err
	}
	return user, nil
}
