package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-booking-system/internal/database/repositories"
	"go-booking-system/internal/models"
	"regexp"
	"testing"
	"time"
)

func TestBookingRepository_Create(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	roomId := 4
	userId := 2
	layout := "2006-01-02 15:04:05"
	dateTimeStartString, dateTimeEndString := "2025-04-03 16:30:00", "2025-04-03 17:00:00"
	dateTimeStart, _ := time.Parse(layout, dateTimeStartString)
	dateTimeEnd, _ := time.Parse(layout, dateTimeEndString)

	bookingToCreate := models.Booking{
		UserId:        userId,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		CreatedBy:     userId,
	}

	rows := sqlmock.NewRows([]string{"booking_id", "user_id", "room_id", "datetime_start", "datetime_end", "active", "created_by", "created_at"}).
		AddRow(1, bookingToCreate.UserId, bookingToCreate.RoomId, bookingToCreate.DateTimeStart, bookingToCreate.DateTimeEnd, true, bookingToCreate.CreatedBy, time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "bookings" ("user_id","room_id","datetime_start","datetime_end","created_by","created_at") VALUES ($1,$2,$3,$4,$5,$6)`,
	)).
		WithArgs(bookingToCreate.UserId, bookingToCreate.RoomId, bookingToCreate.DateTimeStart, bookingToCreate.DateTimeEnd, bookingToCreate.CreatedBy, NotNullTimeArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	// 2. Act
	createdBooking, err := repo.Create(bookingToCreate)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, createdBooking.BookingId)
	assert.Equal(t, 1, createdBooking.BookingId)
	assert.Equal(t, bookingToCreate.UserId, createdBooking.UserId)
	assert.Equal(t, bookingToCreate.RoomId, createdBooking.RoomId)
	assert.Equal(t, bookingToCreate.DateTimeStart, createdBooking.DateTimeStart)
	assert.Equal(t, bookingToCreate.DateTimeEnd, createdBooking.DateTimeEnd)
	assert.Equal(t, true, createdBooking.Active)
	assert.Equal(t, false, createdBooking.CreatedAt.IsZero())
}

func TestBookingRepository_GetAll(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	layout := "2006-01-02 15:04:05"
	dateTimeStart1String, dateTimeEnd1String := "2025-04-03 16:30:00", "2025-04-03 17:00:00"
	dateTimeStart2String, dateTimeEnd2String := "2025-04-03 19:30:00", "2025-04-03 20:00:00"
	dateTimeStart1, _ := time.Parse(layout, dateTimeStart1String)
	dateTimeEnd1, _ := time.Parse(layout, dateTimeEnd1String)
	dateTimeStart2, _ := time.Parse(layout, dateTimeStart2String)
	dateTimeEnd2, _ := time.Parse(layout, dateTimeEnd2String)

	expectedBooking1 := models.Booking{
		UserId:        1,
		RoomId:        2,
		DateTimeStart: dateTimeStart1,
		DateTimeEnd:   dateTimeEnd1,
		CreatedAt:     time.Now(),
	}
	expectedBooking2 := models.Booking{
		UserId:        3,
		RoomId:        4,
		DateTimeStart: dateTimeStart2,
		DateTimeEnd:   dateTimeEnd2,
		CreatedAt:     time.Now(),
	}

	rows := sqlmock.NewRows([]string{"booking_id", "user_id", "room_id", "datetime_start", "datetime_end", "active", "created_at"}).
		AddRow(1, expectedBooking1.UserId, expectedBooking1.RoomId, expectedBooking1.DateTimeStart, expectedBooking1.DateTimeEnd, expectedBooking1.Active, expectedBooking1.CreatedAt).
		AddRow(2, expectedBooking2.UserId, expectedBooking2.RoomId, expectedBooking2.DateTimeStart, expectedBooking2.DateTimeEnd, expectedBooking2.Active, expectedBooking2.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings"`)).
		WillReturnRows(rows)

	// 2. Act
	allBookings := repo.GetAll()

	// 3. Assert
	assert.Equal(t, 1, allBookings[0].BookingId)
	assert.Equal(t, expectedBooking1.UserId, allBookings[0].UserId)
	assert.Equal(t, expectedBooking1.RoomId, allBookings[0].RoomId)
	assert.Equal(t, expectedBooking1.DateTimeStart, allBookings[0].DateTimeStart)
	assert.Equal(t, expectedBooking1.DateTimeEnd, allBookings[0].DateTimeEnd)
	assert.Equal(t, expectedBooking1.Active, allBookings[0].Active)
	assert.Equal(t, expectedBooking1.CreatedAt, allBookings[0].CreatedAt)

	assert.Equal(t, 2, allBookings[1].BookingId)
	assert.Equal(t, expectedBooking2.UserId, allBookings[1].UserId)
	assert.Equal(t, expectedBooking2.RoomId, allBookings[1].RoomId)
	assert.Equal(t, expectedBooking2.DateTimeStart, allBookings[1].DateTimeStart)
	assert.Equal(t, expectedBooking2.DateTimeEnd, allBookings[1].DateTimeEnd)
	assert.Equal(t, expectedBooking2.Active, allBookings[1].Active)
	assert.Equal(t, expectedBooking2.CreatedAt, allBookings[1].CreatedAt)
}

func TestBookingRepository_GetBookingByID(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	layout := "2006-01-02 15:04:05"
	dateTimeStartString, dateTimeEndString := "2025-04-03 16:30:00", "2025-04-03 17:00:00"
	dateTimeStart, _ := time.Parse(layout, dateTimeStartString)
	dateTimeEnd, _ := time.Parse(layout, dateTimeEndString)

	expectedBooking := models.Booking{
		BookingId:     6,
		UserId:        3,
		RoomId:        4,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		Active:        true,
		CreatedAt:     time.Now(),
	}

	rows := sqlmock.NewRows([]string{"booking_id", "user_id", "room_id", "datetime_start", "datetime_end", "active", "created_at"}).
		AddRow(expectedBooking.BookingId, expectedBooking.UserId, expectedBooking.RoomId, expectedBooking.DateTimeStart, expectedBooking.DateTimeEnd, expectedBooking.Active, expectedBooking.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bookings" WHERE "booking_id" = $1`)).
		WithArgs(expectedBooking.BookingId).
		WillReturnRows(rows)

	// 2. Act
	actualBooking, err := repo.GetBookingById(expectedBooking.BookingId)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, actualBooking.BookingId)
	assert.Equal(t, expectedBooking.BookingId, actualBooking.BookingId)
	assert.Equal(t, expectedBooking.UserId, actualBooking.UserId)
	assert.Equal(t, expectedBooking.RoomId, actualBooking.RoomId)
	assert.Equal(t, expectedBooking.DateTimeStart, actualBooking.DateTimeStart)
	assert.Equal(t, expectedBooking.DateTimeEnd, actualBooking.DateTimeEnd)
	assert.Equal(t, true, actualBooking.Active)
	assert.Equal(t, false, actualBooking.CreatedAt.IsZero())
}

func TestBookingRepository_Update(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	layout := "2006-01-02 15:04:05"

	bookingToUpdateBookingId := 30
	newUserId := 999
	newRoomId := 100000
	oldDateTimeEnd, _ := time.Parse(layout, "2025-04-03 16:30:00")
	newDateTimeStart, _ := time.Parse(layout, "2025-04-03 15:30:00")

	updateData := models.Booking{
		BookingId:     bookingToUpdateBookingId,
		UserId:        newUserId,
		RoomId:        newRoomId,
		DateTimeStart: newDateTimeStart,
		DateTimeEnd:   oldDateTimeEnd,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "bookings" SET "user_id"=$1,"room_id"=$2,"datetime_start"=$3,"datetime_end"=$4,"updated_at"=$5  WHERE "booking_id" = $6`,
	)).
		WithArgs(updateData.UserId, updateData.RoomId, NotNullTimeArg(), NotNullTimeArg(), NotNullTimeArg(), updateData.BookingId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	updatedBooking, err := repo.Update(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, updatedBooking.UpdatedAt)
	assert.Equal(t, bookingToUpdateBookingId, updatedBooking.BookingId)
	assert.Equal(t, newUserId, updatedBooking.UserId)
	assert.Equal(t, newRoomId, updatedBooking.RoomId)
	assert.Equal(t, newDateTimeStart, updatedBooking.DateTimeStart)
	assert.Equal(t, oldDateTimeEnd, updatedBooking.DateTimeEnd)
}

func TestBookingRepository_Delete(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	bookingId := 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "bookings" SET "active"=$1,"deleted_at"=$2  WHERE "active"=$3 AND "booking_id" = $4`,
	)).
		WithArgs(false, NotNullTimeArg(), true, bookingId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	isDeleted, err := repo.Delete(bookingId)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, true, isDeleted)
}

func TestBookingRepository_GetBookingsByRoomId(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	roomId := 3
	layout := "2006-01-02 15:04:05"
	dateTimeStart1String, dateTimeEnd1String := "2025-04-03 16:30:00", "2025-04-03 17:00:00"
	dateTimeStart2String, dateTimeEnd2String := "2025-04-03 19:30:00", "2025-04-03 20:00:00"
	dateTimeStart1, _ := time.Parse(layout, dateTimeStart1String)
	dateTimeEnd1, _ := time.Parse(layout, dateTimeEnd1String)
	dateTimeStart2, _ := time.Parse(layout, dateTimeStart2String)
	dateTimeEnd2, _ := time.Parse(layout, dateTimeEnd2String)

	expectedBooking1 := models.Booking{
		UserId:        1,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart1,
		DateTimeEnd:   dateTimeEnd1,
		CreatedAt:     time.Now(),
	}
	expectedBooking2 := models.Booking{
		UserId:        3,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart2,
		DateTimeEnd:   dateTimeEnd2,
		CreatedAt:     time.Now(),
	}

	rows := sqlmock.NewRows([]string{"booking_id", "user_id", "room_id", "datetime_start", "datetime_end", "active", "created_at"}).
		AddRow(1, expectedBooking1.UserId, expectedBooking1.RoomId, expectedBooking1.DateTimeStart, expectedBooking1.DateTimeEnd, expectedBooking1.Active, expectedBooking1.CreatedAt).
		AddRow(2, expectedBooking2.UserId, expectedBooking2.RoomId, expectedBooking2.DateTimeStart, expectedBooking2.DateTimeEnd, expectedBooking2.Active, expectedBooking2.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookings" WHERE "room_id" = $1`,
	)).
		WithArgs(roomId).
		WillReturnRows(rows)

	// 2. Act
	allBookings, err := repo.GetBookingsByRoomId(roomId)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, allBookings[0].BookingId)
	assert.Equal(t, expectedBooking1.UserId, allBookings[0].UserId)
	assert.Equal(t, roomId, allBookings[0].RoomId)
	assert.Equal(t, expectedBooking1.DateTimeStart, allBookings[0].DateTimeStart)
	assert.Equal(t, expectedBooking1.DateTimeEnd, allBookings[0].DateTimeEnd)
	assert.Equal(t, expectedBooking1.Active, allBookings[0].Active)
	assert.Equal(t, expectedBooking1.CreatedAt, allBookings[0].CreatedAt)

	assert.Equal(t, 2, allBookings[1].BookingId)
	assert.Equal(t, expectedBooking2.UserId, allBookings[1].UserId)
	assert.Equal(t, roomId, allBookings[1].RoomId)
	assert.Equal(t, expectedBooking2.DateTimeStart, allBookings[1].DateTimeStart)
	assert.Equal(t, expectedBooking2.DateTimeEnd, allBookings[1].DateTimeEnd)
	assert.Equal(t, expectedBooking2.Active, allBookings[1].Active)
	assert.Equal(t, expectedBooking2.CreatedAt, allBookings[1].CreatedAt)
}

func TestBookingRepository_GetBookingsByRoomIdAndBookingTime(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewBookingRepositoryPostgres(db)

	roomId := 3

	layout := "2006-01-02 15:04:05"
	timeStartString, timeEndString := "2025-04-03 16:00:00", "2025-04-03 21:00:00"
	dateTimeStart1String, dateTimeEnd1String := "2025-04-03 16:30:00", "2025-04-03 17:00:00"
	dateTimeStart2String, dateTimeEnd2String := "2025-04-03 19:30:00", "2025-04-03 20:00:00"
	dateTimeStart1, _ := time.Parse(layout, dateTimeStart1String)
	dateTimeEnd1, _ := time.Parse(layout, dateTimeEnd1String)
	dateTimeStart2, _ := time.Parse(layout, dateTimeStart2String)
	dateTimeEnd2, _ := time.Parse(layout, dateTimeEnd2String)
	dateTimeStart, _ := time.Parse(layout, timeStartString)
	dateTimeEnd, _ := time.Parse(layout, timeEndString)

	expectedBooking1 := models.Booking{
		UserId:        1,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart1,
		DateTimeEnd:   dateTimeEnd1,
		CreatedAt:     time.Now(),
	}
	expectedBooking2 := models.Booking{
		UserId:        3,
		RoomId:        roomId,
		DateTimeStart: dateTimeStart2,
		DateTimeEnd:   dateTimeEnd2,
		CreatedAt:     time.Now(),
	}

	rows := sqlmock.NewRows([]string{"booking_id", "user_id", "room_id", "datetime_start", "datetime_end", "active", "created_at"}).
		AddRow(1, expectedBooking1.UserId, expectedBooking1.RoomId, expectedBooking1.DateTimeStart, expectedBooking1.DateTimeEnd, expectedBooking1.Active, expectedBooking1.CreatedAt).
		AddRow(2, expectedBooking2.UserId, expectedBooking2.RoomId, expectedBooking2.DateTimeStart, expectedBooking2.DateTimeEnd, expectedBooking2.Active, expectedBooking2.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "bookings" WHERE room_id = $1 `+
			`AND `+
			`(( `+
			`($2 BETWEEN datetime_start AND datetime_end) `+
			`OR ($3 BETWEEN datetime_start AND datetime_end) `+
			`OR (datetime_start BETWEEN $4 AND $5) `+
			`OR (datetime_end BETWEEN $6 AND $7)`+
			`))`,
	)).
		WithArgs(roomId, NotNullTimeArg(), NotNullTimeArg(), NotNullTimeArg(), NotNullTimeArg(), NotNullTimeArg(), NotNullTimeArg()).
		WillReturnRows(rows)

	// 2. Act
	overlappingBookings, err := repo.GetBookingsByRoomIdAndBookingTime(roomId, dateTimeStart, dateTimeEnd)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, overlappingBookings[0].BookingId)
	assert.Equal(t, expectedBooking1.UserId, overlappingBookings[0].UserId)
	assert.Equal(t, roomId, overlappingBookings[0].RoomId)
	assert.Equal(t, expectedBooking1.DateTimeStart, overlappingBookings[0].DateTimeStart)
	assert.Equal(t, expectedBooking1.DateTimeEnd, overlappingBookings[0].DateTimeEnd)
	assert.Equal(t, expectedBooking1.Active, overlappingBookings[0].Active)
	assert.Equal(t, expectedBooking1.CreatedAt, overlappingBookings[0].CreatedAt)

	assert.Equal(t, 2, overlappingBookings[1].BookingId)
	assert.Equal(t, expectedBooking2.UserId, overlappingBookings[1].UserId)
	assert.Equal(t, roomId, overlappingBookings[1].RoomId)
	assert.Equal(t, expectedBooking2.DateTimeStart, overlappingBookings[1].DateTimeStart)
	assert.Equal(t, expectedBooking2.DateTimeEnd, overlappingBookings[1].DateTimeEnd)
	assert.Equal(t, expectedBooking2.Active, overlappingBookings[1].Active)
	assert.Equal(t, expectedBooking2.CreatedAt, overlappingBookings[1].CreatedAt)

}
