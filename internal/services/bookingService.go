package services

import (
	"fmt"
	"go-booking-system/internal/database"
	"go-booking-system/internal/models"
	"sort"
	"time"
)

const (
	RoomIsNotAvailable = false
	RoomIsAvailable    = true
	NotOverlapping     = false
	Overlapping        = true
)

type BookingService struct {
	repository database.BookingRepository
}

func NewBookingService(repository database.BookingRepository) *BookingService {
	return &BookingService{repository: repository}
}

func (b *BookingService) GetAll() []models.Booking {
	return b.repository.GetAll()
}

func (b *BookingService) GetBookingById(bookingId int) (models.Booking, error) {
	return b.repository.GetBookingById(bookingId)
}

func (b *BookingService) GetBookingsByRoomId(roomId int) ([]models.Booking, error) {
	return b.repository.GetBookingsByRoomId(roomId)
}

func (b *BookingService) GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error) {
	return b.repository.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
}

func (b *BookingService) Update(booking models.Booking) (models.Booking, error) {
	// 1. Check if time slot [start; end] is available
	isRoomAvailable, err := b.CheckIfRoomAvailable(booking.RoomId, booking.DateTimeStart, booking.DateTimeEnd)
	if err != nil {
		return models.Booking{}, err
	}

	// 2. If time slot [start; end] is available => update booking
	if isRoomAvailable {
		return b.repository.Update(booking)
	}

	// 3. check overlapping bookings
	overlappingBookings, err := b.repository.GetBookingsByRoomIdAndBookingTime(
		booking.RoomId, booking.DateTimeStart, booking.DateTimeEnd)
	if err != nil {
		return models.Booking{}, err
	}

	// 4. if only one overlapping booking AND it has same id as passed booking => we can update booking
	if len(overlappingBookings) == 1 && overlappingBookings[0].BookingId == booking.BookingId {
		return models.Booking{}, err
	}

	return models.Booking{}, NewOverlappingBookingsError("cannot update booking. Overlapping bookings exist.", overlappingBookings)
}

func (b *BookingService) Delete(bookingId int) (bool, error) {
	return b.repository.Delete(bookingId)
}

// CheckIfRoomAvailable true - available; false - not available
func (b *BookingService) CheckIfRoomAvailable(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) (bool, error) {
	bookingToCheck := models.Booking{
		RoomId:        roomId,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
	}

	// 1. get overlapping bookings
	overlappingBookings, err := b.GetOverlappingBookings(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		return false, err
	}

	// 2. Check if our Booking is Overlapping with other bookings
	isOverlapping, err := b.IsOverlapping(bookingToCheck, overlappingBookings...)
	if err != nil {
		return false, err
	}

	// 3. if room is available during [start; end]
	if isOverlapping == false {
		return RoomIsAvailable, nil
	} else {
		return RoomIsNotAvailable, NewOverlappingBookingsError("Room is not available", overlappingBookings)
	}
}

func (b *BookingService) BookRoom(userId int, roomId int, dateTimeStart time.Time, dateTimeEnd time.Time, createdBy int) (models.Booking, error) {
	bookingToCreate := models.Booking{
		UserId:        userId,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		CreatedBy:     createdBy,
	}

	// 1. Check if it is possible to book Room in the given timeframe [start; end]
	available, err := b.CheckIfRoomAvailable(roomId, dateTimeStart, dateTimeEnd)
	if err != nil {
		return models.Booking{}, err
	}
	// 2. if room is available during [start; end]
	if available {
		return b.repository.Create(bookingToCreate)
	}

	return models.Booking{}, fmt.Errorf("room booking ended with an error. Passed data: %v", bookingToCreate)
}

func (b *BookingService) GetOverlappingBookings(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error) {
	return b.repository.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)
}

func (b *BookingService) IsOverlapping(bookingToCheck models.Booking, overlapingBookings ...models.Booking) (bool, error) {
	// 1.1 If no overlapping bookings found => success
	if len(overlapingBookings) == 0 {
		return false, nil
	} else
	// 1.2 Check if booking ends at start of another || no overlapping
	if len(overlapingBookings) == 1 {
		return b.CheckOverlapping(bookingToCheck, overlapingBookings[0]), nil
	} else
	// 1.3 Check if booking is between two bookings
	if len(overlapingBookings) == 2 {
		return b.CheckOverlapping(overlapingBookings[0], bookingToCheck) ||
			b.CheckOverlapping(bookingToCheck, overlapingBookings[1]), nil
	} else {
		return true, nil
	}
}

func (b *BookingService) CheckOverlapping(b1 models.Booking, b2 models.Booking) bool {
	// 1. Sort passed bookings
	sortedBookings := b.SortBookings(b1, b2)
	b1, b2 = sortedBookings[0], sortedBookings[1]

	// 2.1 Supplementary expressions for logical expression: "Is booking #1 overlapping with booking #2?"
	B1StartsBeforeB2Starts := b1.DateTimeStart.Before(b2.DateTimeStart)
	B1endEqualsB2start := b1.DateTimeEnd == b2.DateTimeStart
	B1EndsBeforeB2Ends := b1.DateTimeEnd.Before(b2.DateTimeEnd)
	B1EndsBeforeB2Starts := b1.DateTimeEnd.Before(b2.DateTimeStart)
	B1EndsBeforeB2StartsOrEquals := B1EndsBeforeB2Starts || B1endEqualsB2start

	// 2.2 Logical expression: "Is booking #1 overlapping with booking #2?"
	isNotOverlapping := B1StartsBeforeB2Starts && B1EndsBeforeB2Ends && B1EndsBeforeB2StartsOrEquals

	// 3.1 if booking #1 doesn't overlap booking #2
	if isNotOverlapping {
		return NotOverlapping
	}

	// 3.2 if booking #1 overlaps booking #2
	return Overlapping
}

func (b *BookingService) SortBookings(bookings ...models.Booking) []models.Booking {
	sort.Slice(bookings, func(i, j int) bool {
		b1 := bookings[i]
		b2 := bookings[j]
		return (b1.DateTimeStart.Before(b2.DateTimeStart)) ||
			((b1.DateTimeStart == b2.DateTimeStart) && (b1.DateTimeEnd.Before(b2.DateTimeEnd)))
	})
	return bookings
}

type OverlappingBookingsError struct {
	Message             string
	OverlappingBookings []models.Booking
}

func NewOverlappingBookingsError(message string, overlappingBookings []models.Booking) *OverlappingBookingsError {
	return &OverlappingBookingsError{Message: message, OverlappingBookings: overlappingBookings}
}

func (o *OverlappingBookingsError) Error() string {
	return fmt.Sprintf(`{message: '%v', overlapping_bookings: '%v'}`, o.Message, o.OverlappingBookings)
}
