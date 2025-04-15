package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
	"log"
	"time"
)

type UserRepository struct {
	connection *gorm.DB
}

func (u *UserRepository) Create(user models.User) (models.User, error) {
	result := u.connection.
		Omit("updated_at", "deleted_at").
		Select("name", "email", "telephone", "role_id", "username", "password", "active").
		Create(&user)

	if err := result.Error; err != nil {
		log.Println("UserRepository.Create(): error occured during User creation. Passed data: ", user)
		log.Println(err)
		return user, err
	}

	return user, nil
}

func (u *UserRepository) GetAll() []models.User {
	var allUsers []models.User

	u.connection.Find(&allUsers)

	return allUsers
}

func (u *UserRepository) GetUserById(userId int) (models.User, error) {
	var foundUser models.User
	result := u.connection.Find(&foundUser, "user_id", userId)

	if err := result.Error; err != nil {
		log.Println("UserRepository.GetUserById(): error occured during User search. Passed data: ", userId)
		log.Println(err)
		return models.User{}, err
	}

	return foundUser, nil
}

func (u *UserRepository) Update(user models.User) (models.User, error) {
	result := u.connection.
		Omit("active", "created_at", "deleted_at"). // `active` is changed only at DELETION
		Model(&user).
		Updates(&user)

	if err := result.Error; err != nil {
		log.Println("UserRepository.Update(): error occured during User update. Passed data: ", user)
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) Delete(userId int) (bool, error) {
	userToDelete := models.User{
		UserId:    userId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := u.connection.
		Select("*").
		Omit("created_at", "updated_at", "role_id", "name", "email", "telephone", "username", "password").
		Model(&userToDelete).
		Updates(&userToDelete)

	if err := result.Error; err != nil {
		log.Println("UserRepository.Delete(): error occured during User deletion. Passed data: ", userId)
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (u *UserRepository) UpdatePassword(user models.User) (models.User, error) {
	result := u.connection.
		Omit("name", "email", "telephone", "role_id", "username", "active", "created_at", "deleted_at").
		Model(&user).
		Updates(&user)

	if err := result.Error; err != nil {
		log.Println("UserRepository.UpdatePassword(): error occured during password change. Passed data: ", user.UserId, user.Password)
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) UpdateUsername(user models.User) (models.User, error) {
	result := u.connection.
		Omit("name", "email", "telephone", "role_id", "password", "active", "created_at", "deleted_at").
		Model(&user).
		Updates(&user)

	if err := result.Error; err != nil {
		log.Println("UserRepository.UpdateUsername(): error occured during username change. Passed data: ", user.UserId, user.UserName)
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) UpdateUserRole(user models.User) (models.User, error) {
	result := u.connection.
		Omit("name", "email", "telephone", "username", "password", "active", "created_at", "deleted_at").
		Model(&user).
		Updates(&user)

	if err := result.Error; err != nil {
		log.Println("UserRepository.UpdateUserRole(): error occured during user's role change. Passed data: ", user.UserId, user.RoleId)
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func NewUserRepositoryPostgres(connection *gorm.DB) *UserRepository {
	return &UserRepository{connection: connection}
}
