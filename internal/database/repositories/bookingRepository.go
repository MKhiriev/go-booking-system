package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
	"time"
)

type BookingRepository struct {
	connection *gorm.DB
}

func (b *BookingRepository) Update(booking models.Booking) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) Delete(bookingId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) CheckIfRoomAvailable(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, bool) {
	//TODO implement me
	/*
		SELECT *
		FROM Booking
		WHERE room_id = @room_id
		AND (
		    (@datetime_start BETWEEN datetime_start AND datetime_end)
		    OR (@datetime_end BETWEEN datetime_start AND datetime_end)
		    OR (datetime_start BETWEEN @datetime_start AND @datetime_end)
		    OR (datetime_end BETWEEN @datetime_start AND @datetime_end)
		);
	*/
	panic("implement me")
}

func NewBookingRepositoryPostgres(connection *gorm.DB) *BookingRepository {
	return &BookingRepository{connection: connection}
}

func (b *BookingRepository) Create(booking models.Booking) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) GetAll() []models.Booking {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) GetBookingById(bookingId int) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) GetBookingsByRoomId(roomId int) ([]models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookingRepository) BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error) {
	//TODO implement me
	panic("implement me")
}
