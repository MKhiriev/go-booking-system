package database

import (
	"gorm.io/gorm"
	"humoBooking/internal/database/repositories"
	"humoBooking/internal/models"
	"time"
)

type Database struct {
	BookingRepository
	RoomRepository
	UserRepository
}

func NewDatabase(conn *gorm.DB) *Database {
	return &Database{
		BookingRepository: repositories.NewBookingRepositoryPostgres(conn),
		RoomRepository:    repositories.NewRoomRepositoryPostgres(conn),
		UserRepository:    repositories.NewUserRepositoryPostgres(conn),
	}
}

type BookingRepository interface {
	Create(booking models.Booking) (models.Booking, error)
	GetAll() []models.Booking
	GetBookingById(bookingId int) (models.Booking, error)
	GetBookingsByRoomId(roomId int) ([]models.Booking, error)
	GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error)
	Update(booking models.Booking) (models.Booking, error)
	Delete(bookingId int) (bool, error)
}

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetAll() []models.User
	GetUserById(userId int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(userId int) (bool, error)
	UpdatePassword(user models.User) (models.User, error)
	UpdateUsername(user models.User) (models.User, error)
	UpdateUserRole(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type RoomRepository interface {
	Create(room models.Room) (models.Room, error)
	GetAll() []models.Room
	GetRoomById(roomId int) (models.Room, error)
	Update(room models.Room) (models.Room, error)
	Delete(roomId int) (bool, error)
}
