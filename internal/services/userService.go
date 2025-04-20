package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
)

type UserService struct {
	repository database.UserRepository
}

func NewUserService(repository database.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (u *UserService) GetAll() []models.User {
	return u.repository.GetAll()
}

func (u *UserService) GetUserById(userId int) (models.User, error) {
	return u.repository.GetUserById(userId)
}

func (u *UserService) Update(user models.User) (models.User, error) {
	return u.repository.Update(user)
}

func (u *UserService) Delete(userId int) (bool, error) {
	return u.repository.Delete(userId)
}
