package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
	"time"
)

type Service struct {
	BookingService BookingServiceInterface
	RoomService    RoomServiceInterface
	UserService    UserServiceInterface
}

func NewService(db *database.Database) *Service {
	return &Service{
		BookingService: NewBookingService(db.BookingRepository),
		RoomService:    NewRoomService(db.RoomRepository),
		UserService:    NewUserService(db.UserRepository),
	}
}

type BookingServiceInterface interface {
	CheckIfRoomAvailable(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error)
	BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error)
	GetAll() []models.Booking
	GetBookingById(bookingId int) (models.Booking, error)
	GetBookingsByRoomId(roomId int) ([]models.Booking, error)
	GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error)
	Update(booking models.Booking) (models.Booking, error)
	Delete(bookingId int) (bool, error)
}

type RoomServiceInterface interface {
	GetAllRooms() []models.Room
	Create(room models.Room) (models.Room, error)
	GetAll() []models.Room
	GetRoomById(roomId int) (models.Room, error)
	Update(room models.Room) (models.Room, error)
	Delete(roomId int) (bool, error)
}

type UserServiceInterface interface {
	Create(user models.User) (models.Room, error)
	GetAll() []models.User
	GetUserById(userId int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(userId int) (bool, error)
}
