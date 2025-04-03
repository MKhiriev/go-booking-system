package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
)

type UserService struct {
	database.UserRepository
}

func (u *UserService) Create(user models.User) (models.Room, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) GetAll() []models.User {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) GetUserById(userId int) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) Update(user models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) Delete(userId int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService(repository database.UserRepository) *UserService {
	return &UserService{UserRepository: repository}
}
