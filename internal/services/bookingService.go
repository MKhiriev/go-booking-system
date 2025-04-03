package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
	"time"
)

type BookingService struct {
	database.BookingRepository
}

func NewBookingService(repository database.BookingRepository) *BookingService {
	return &BookingService{BookingRepository: repository}
}

func (b *BookingService) GetAll() []models.Booking {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) GetBookingById(bookingId int) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) GetBookingsByRoomId(roomId int) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) Update(booking models.Booking) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) Delete(bookingId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) CheckIfRoomAvailable(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingService) BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}
