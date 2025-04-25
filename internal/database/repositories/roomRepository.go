package repositories

import (
	"errors"
	"gorm.io/gorm"
	"humoBooking/internal/models"
	"log"
	"time"
)

type RoomRepository struct {
	connection *gorm.DB
}

func NewRoomRepositoryPostgres(connection *gorm.DB) *RoomRepository {
	return &RoomRepository{connection: connection}
}

func (r *RoomRepository) Create(room models.Room) (models.Room, error) {
	result := r.connection.
		Omit("room_id", "updated_at", "deleted_at").
		Select("number", "capacity", "created_by").
		Create(&room)

	if err := result.Error; err != nil {
		log.Println("RoomRepository.Create(): error occured during Room creation. Passed data: ", room)
		log.Println(err)
		return models.Room{}, err
	}

	return room, nil
}

func (r *RoomRepository) GetAll() []models.Room {
	var allRooms []models.Room

	r.connection.Find(&allRooms)

	return allRooms
}

func (r *RoomRepository) GetRoomById(roomId int) (models.Room, error) {
	var foundRoom models.Room

	result := r.connection.Find(&foundRoom, "room_id", roomId)
	if err := result.Error; err != nil {
		log.Println("RoomRepository.GetRoomById(): error occured during Room search. Passed data: ", roomId)
		log.Println(err)
		return models.Room{}, err
	}

	if rowsReturned := result.RowsAffected; rowsReturned == 0 {
		log.Println("RoomRepository.GetRoomById(): no Rooms were found. Passed data: ", roomId)
		return models.Room{}, errors.New("no Rooms were found")
	}

	return foundRoom, nil
}

func (r *RoomRepository) Update(room models.Room) (models.Room, error) {
	result := r.connection.
		Omit("active", "created_at", "deleted_at").
		Model(&room).
		Updates(&room)

	if err := result.Error; err != nil {
		log.Println("RoomRepository.Update(): error occured during Room update. Passed data: ", room)
		return room, err
	}

	if rowsUpdated := result.RowsAffected; rowsUpdated == 0 {
		log.Println("RoomRepository.Update(): no Rooms were updated. Reason: Room to update not found. Passed data: ", room)
		return room, errors.New("no Rooms were updated")
	}

	return room, nil
}

func (r *RoomRepository) Delete(roomId int) (bool, error) {
	roomToDelete := models.Room{
		RoomId:    roomId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := r.connection.
		Select("*").
		Where("active = true").
		Omit("number", "capacity", "created_by", "created_at", "updated_at").
		Model(&roomToDelete).
		Updates(&roomToDelete)

	if err := result.Error; err != nil {
		log.Println("RoomRepository.Delete(): error occured during Room deletion. Passed data: ", roomId)
		log.Println(err)
		return false, err
	}

	if rowsDeleted := result.RowsAffected; rowsDeleted == 0 {
		log.Println("RoomRepository.Delete(): no Rooms were deleted. Reason: Room to delete not found. Passed data: ", roomId)
		return false, errors.New("no Rooms were deleted")
	}

	return true, nil
}
