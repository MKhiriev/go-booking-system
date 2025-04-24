package repositories

import (
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
		Omit("updated_at", "deleted_at").
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
		Omit("created_at", "updated_at", "number", "capacity").
		Model(&roomToDelete).
		Updates(&roomToDelete)

	if err := result.Error; err != nil {
		log.Println("RoomRepository.Delete(): error occured during Room deletion. Passed data: ", roomId)
		log.Println(err)
		return false, err
	}

	return true, nil
}
