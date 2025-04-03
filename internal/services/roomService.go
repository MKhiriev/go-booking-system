package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
)

type RoomService struct {
	database.RoomRepository
}

func NewRoomService(repository database.RoomRepository) *RoomService {
	return &RoomService{RoomRepository: repository}
}

func (r *RoomService) GetAllRooms() []models.Room {

	//TODO implement me
	panic("implement me")
}

func (r *RoomService) Create(room models.Room) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) GetAll() []models.Room {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) GetRoomById(roomId int) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) Update(room models.Room) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomService) Delete(roomId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}
