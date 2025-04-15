package models

import (
	"fmt"
	"time"
)

type Booking struct {
	BookingId     int       `json:"booking_id" gorm:"primarykey"`
	UserId        int       `json:"user_id"`
	RoomId        int       `json:"room_id"`
	DateTimeStart time.Time `json:"datetime_start" gorm:"column:datetime_start"`
	DateTimeEnd   time.Time `json:"datetime_end" gorm:"column:datetime_end"`

	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (b Booking) String() string {
	return fmt.Sprintf("Booking {id: %d | time: %s - %s | date: %s | room_id: %d | user_id: %d}", b.BookingId, b.DateTimeStart.Format("15:04:05"), b.DateTimeEnd.Format("15:04:05"), b.DateTimeStart.Format("2006-01-02"), b.RoomId, b.UserId)
}
