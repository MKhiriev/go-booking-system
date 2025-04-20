package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"time"
)

type Service struct {
	BookingService BookingServiceInterface
	RoomService    RoomServiceInterface
	UserService    UserServiceInterface
	AuthService    AuthServiceInterface
}

func NewService(db *database.Database) *Service {
	return &Service{
		BookingService: NewBookingService(db.BookingRepository),
		RoomService:    NewRoomService(db.RoomRepository),
		UserService:    NewUserService(db.UserRepository),
		AuthService:    NewAuthService(db.UserRepository),
	}
}

type AuthServiceInterface interface {
	Create(user models.User) (models.User, error)
	UpdatePassword(userId int, password string) (models.User, error)
	UpdateUsername(userId int, username string) (models.User, error)
	UpdateRole(userId int, roleId int) (models.User, error)
	CheckIfUserExistsAndPasswordIsCorrect(username string, password string) (models.User, error)
	GeneratePasswordHash(password string) string
	GenerateTokens(user models.User, identity pkg.IPAddressIdentity) (accessToken pkg.JWTToken, refreshToken pkg.JWTToken)
	ValidateAccessToken(encodedToken string, ipAddress string) *JWTTokenValidator
	ValidateRefreshToken(encodedToken string, ipAddress string) *JWTTokenValidator
}

type BookingServiceInterface interface {
	CheckIfRoomAvailable(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (bool, error)
	BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error)
	GetAll() []models.Booking
	GetBookingById(bookingId int) (models.Booking, error)
	GetBookingsByRoomId(roomId int) ([]models.Booking, error)
	GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error)
	Update(booking models.Booking) (models.Booking, error)
	Delete(bookingId int) (bool, error)
}

type RoomServiceInterface interface {
	Create(room models.Room) (models.Room, error)
	GetAll() []models.Room
	GetRoomById(roomId int) (models.Room, error)
	Update(room models.Room) (models.Room, error)
	Delete(roomId int) (bool, error)
}

type UserServiceInterface interface {
	Create(user models.User) (models.User, error)
	GetAll() []models.User
	GetUserById(userId int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(userId int) (bool, error)
}
