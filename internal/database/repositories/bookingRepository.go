package repositories

import (
	"database/sql"
	"errors"
	"go-booking-system/internal/models"
	"gorm.io/gorm"
	"log"
	"time"
)

// BookingRepository TODO create custom errors
type BookingRepository struct {
	connection *gorm.DB
}

func (b *BookingRepository) Update(booking models.Booking) (models.Booking, error) {
	result := b.connection.
		Omit("active", "created_at", "deleted_at"). // `active` is changed only at DELETION
		Model(&booking).
		Updates(&booking)

	if err := result.Error; err != nil {
		log.Println("BookingRepository.Update(): error occured during Booking update. Passed data: ", booking)
		log.Println(err)
		return models.Booking{}, err
	}

	if rowsUpdated := result.RowsAffected; rowsUpdated == 0 {
		log.Println("BookingRepository.Update(): no Bookings were updated. Reason: Booking to update not found. Passed data: ", booking)
		return booking, errors.New("no Bookings were updated")
	}

	return booking, nil
}

func (b *BookingRepository) Delete(bookingId int) (bool, error) {
	bookingToDelete := models.Booking{
		BookingId: bookingId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := b.connection.
		Select("*").
		Where(`"active"=?`, true).
		Omit("created_by", "created_at", "updated_at", "room_id", "user_id", "datetime_start", "datetime_end").
		Model(&bookingToDelete).
		Updates(&bookingToDelete)

	if err := result.Error; err != nil {
		log.Println("BookingRepository.Delete(): error occured during Booking deletion. Passed data: ", bookingId)
		log.Println(err)
		return false, err
	}

	if rowsDeleted := result.RowsAffected; rowsDeleted == 0 {
		log.Println("BookingRepository.Delete(): no Bookings were deleted. Reason: Booking to delete not found. Passed data: ", bookingId)
		return false, errors.New("no Bookings were deleted")
	}

	return true, nil
}

func NewBookingRepositoryPostgres(connection *gorm.DB) *BookingRepository {
	return &BookingRepository{connection: connection}
}

func (b *BookingRepository) Create(booking models.Booking) (models.Booking, error) {
	result := b.connection.
		Omit("updated_at", "deleted_at", "active").
		Create(&booking)

	if err := result.Error; err != nil {
		log.Println("BookingRepository.Create(): error occured during Booking creation. Passed data: ", booking)
		log.Println(err)
		return models.Booking{}, err
	}

	return booking, nil
}

func (b *BookingRepository) GetAll() []models.Booking {
	var allBookings []models.Booking

	b.connection.Find(&allBookings)

	return allBookings
}

func (b *BookingRepository) GetBookingById(bookingId int) (models.Booking, error) {
	var foundBooking models.Booking

	result := b.connection.Find(&foundBooking, "booking_id", bookingId)
	if err := result.Error; err != nil {
		log.Println("BookingRepository.GetBookingById(): error occured during Booking search. Passed data: ", bookingId)
		log.Println(err)
		return models.Booking{}, err
	}

	if rowsReturned := result.RowsAffected; rowsReturned == 0 {
		log.Println("BookingRepository.GetBookingById(): no Rooms were found. Passed data: ", bookingId)
		return models.Booking{}, errors.New("no Rooms were found")
	}

	return foundBooking, nil
}

func (b *BookingRepository) GetBookingsByRoomId(roomId int) ([]models.Booking, error) {
	var foundBookings []models.Booking

	result := b.connection.Find(&foundBookings, "room_id", roomId)
	if err := result.Error; err != nil {
		log.Println("BookingRepository.GetBookingsByRoomId(): error occured during Bookings search by RoomId. Passed data: ", roomId)
		log.Println(err)
		return nil, err
	}

	return foundBookings, nil
}

func (b *BookingRepository) GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error) {
	var overlapingBookings []models.Booking

	result := b.connection.
		Where("room_id = @room_id", sql.Named("room_id", roomId)).
		Where(`(
		    (@datetime_start BETWEEN datetime_start AND datetime_end)
		    OR (@datetime_end BETWEEN datetime_start AND datetime_end)
		    OR (datetime_start BETWEEN @datetime_start AND @datetime_end)
		    OR (datetime_end BETWEEN @datetime_start AND @datetime_end))`,
			sql.Named("datetime_start", dateTimeStart),
			sql.Named("datetime_end", dateTimeEnd)).
		Find(&overlapingBookings)

	if err := result.Error; err != nil {
		log.Println("BookingRepository.GetBookingsByRoomIdAndBookingTime(): error occured during overlaping Booking search. Passed data: ", roomId, dateTimeStart, dateTimeEnd)
		log.Println(err)
		return nil, err
	}

	return overlapingBookings, nil
}

// BookRoom Probably Service level...
func (b *BookingRepository) BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (models.Booking, error) {
	requestedBooking := models.Booking{
		UserId:        userId,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
	}

	booking, err := b.Create(requestedBooking)
	if err != nil {
		log.Println("BookingRepository.BookRoom(): error occured during Bookinging process. Passed data: ", booking)
		log.Println(err)
		return models.Booking{}, err
	}

	return booking, nil
}
