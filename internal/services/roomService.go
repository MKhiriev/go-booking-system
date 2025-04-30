package services

import (
	"go-booking-system/internal/database"
	"go-booking-system/internal/models"
)

type RoomService struct {
	repository database.RoomRepository
}

func NewRoomService(repository database.RoomRepository) *RoomService {
	return &RoomService{repository: repository}
}

func (r *RoomService) Create(room models.Room) (models.Room, error) {
	return r.repository.Create(room)
}

func (r *RoomService) GetAll() []models.Room {
	return r.repository.GetAll()
}

func (r *RoomService) GetRoomById(roomId int) (models.Room, error) {
	return r.repository.GetRoomById(roomId)
}

func (r *RoomService) Update(room models.Room) (models.Room, error) {
	return r.repository.Update(room)
}

func (r *RoomService) Delete(roomId int) (bool, error) {
	return r.repository.Delete(roomId)
}
