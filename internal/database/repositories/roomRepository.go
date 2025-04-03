package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
)

type RoomRepository struct {
	connection *gorm.DB
}

func NewRoomRepositoryPostgres(connection *gorm.DB) *RoomRepository {
	return &RoomRepository{connection: connection}
}

func (r *RoomRepository) Create(room models.Room) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepository) GetAll() []models.Room {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepository) GetRoomById(roomId int) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepository) Update(room models.Room) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoomRepository) Delete(roomId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}
