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

func TestRoomRepository_Create(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewRoomRepositoryPostgres(db)

	newRoomId := 1
	roomToCreate := models.Room{
		Number:    "Briefing Room #1",
		Capacity:  20,
		CreatedBy: 1,
	}

	rows := sqlmock.NewRows([]string{"room_id", "number", "capacity", "active", "created_at"}).
		AddRow(newRoomId, roomToCreate.Number, roomToCreate.Capacity, true, time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "rooms" ("number","capacity","created_by","created_at") VALUES ($1,$2,$3,$4)`,
	)).
		WithArgs(roomToCreate.Number, roomToCreate.Capacity, roomToCreate.CreatedBy, NotNullTimeArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	// 2. Act
	createdRoom, err := repo.Create(roomToCreate)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, createdRoom.RoomId)
	assert.Equal(t, newRoomId, createdRoom.RoomId)
	assert.Equal(t, roomToCreate.Number, createdRoom.Number)
	assert.Equal(t, roomToCreate.Capacity, createdRoom.Capacity)
	assert.Equal(t, true, createdRoom.Active)
	assert.Equal(t, false, createdRoom.CreatedAt.IsZero())
}

func TestRoomRepository_GetAll(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewRoomRepositoryPostgres(db)

	expectedRoom1 := models.Room{
		Number:   "Conference Room #1",
		Capacity: 21,
		Active:   true,
	}
	expectedRoom2 := models.Room{
		Number:   "Conference Room #2",
		Capacity: 21,
		Active:   true,
	}

	rows := sqlmock.NewRows([]string{"room_id", "number", "capacity", "active"}).
		AddRow(expectedRoom1.RoomId, expectedRoom1.Number, expectedRoom1.Capacity, expectedRoom1.Active).
		AddRow(expectedRoom2.RoomId, expectedRoom2.Number, expectedRoom2.Capacity, expectedRoom2.Active)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WillReturnRows(rows)

	// 2. Act
	allRooms := repo.GetAll()

	// 3. Assert
	assert.Equal(t, expectedRoom1.RoomId, allRooms[0].RoomId)
	assert.Equal(t, expectedRoom1.Number, allRooms[0].Number)
	assert.Equal(t, expectedRoom1.Active, allRooms[0].Active)

	assert.Equal(t, expectedRoom2.RoomId, allRooms[1].RoomId)
	assert.Equal(t, expectedRoom2.Number, allRooms[1].Number)
	assert.Equal(t, expectedRoom2.Active, allRooms[1].Active)
}

func TestRoomRepository_GetUserByID(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewRoomRepositoryPostgres(db)

	expectedRoom := models.Room{
		RoomId:   2,
		Number:   "Conference Room #2",
		Capacity: 21,
		Active:   true,
	}

	rows := sqlmock.NewRows([]string{"room_id", "number", "capacity", "active"}).
		AddRow(expectedRoom.RoomId, expectedRoom.Number, expectedRoom.Capacity, expectedRoom.Active)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "room_id" = $1`)).
		WithArgs(expectedRoom.RoomId).
		WillReturnRows(rows)

	// 2. Act
	actualRoom, err := repo.GetRoomById(expectedRoom.RoomId)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRoom.RoomId, actualRoom.RoomId)
	assert.Equal(t, expectedRoom.Number, actualRoom.Number)
	assert.Equal(t, expectedRoom.Capacity, actualRoom.Capacity)
	assert.Equal(t, expectedRoom.Active, actualRoom.Active)
}

func TestRoomRepository_Update(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewRoomRepositoryPostgres(db)

	room := models.Room{
		RoomId:    1,
		Number:    "Conference Room #1",
		Capacity:  10,
		Active:    true,
		CreatedAt: time.Now(),
	}

	newNumber := "Conference Room #1 - UPDATED"
	newCapacity := 100
	updateData := models.Room{
		RoomId:   room.RoomId,
		Number:   newNumber,
		Capacity: newCapacity,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "rooms" SET "number"=$1,"capacity"=$2,"updated_at"=$3  WHERE "room_id" = $4`,
	)).
		WithArgs(updateData.Number, updateData.Capacity, NotNullTimeArg(), updateData.RoomId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	updatedRoom, err := repo.Update(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, updatedRoom.UpdatedAt)
	assert.Equal(t, room.RoomId, updatedRoom.RoomId)
	assert.Equal(t, newNumber, updatedRoom.Number)
	assert.Equal(t, newCapacity, updatedRoom.Capacity)
}

func TestRoomRepository_Delete(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewRoomRepositoryPostgres(db)

	roomToDelete := models.Room{
		RoomId:    1,
		Number:    "Conference Room #1",
		Capacity:  10,
		Active:    true,
		CreatedAt: time.Now(),
	}
	roomId := roomToDelete.RoomId

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "rooms" SET "active"=$1,"deleted_at"=$2 WHERE "active"=$3 AND "room_id" = $4`,
	)).
		WithArgs(false, NotNullTimeArg(), true, roomId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	isDeleted, err := repo.Delete(roomId)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, true, isDeleted)
}
