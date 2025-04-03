package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
)

type UserRepository struct {
	connection *gorm.DB
}

func (u *UserRepository) Create(user models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetAll() []models.User {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetUserById(userId int) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Update(user models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Delete(userId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepositoryPostgres(connection *gorm.DB) *UserRepository {
	return &UserRepository{connection: connection}
}
