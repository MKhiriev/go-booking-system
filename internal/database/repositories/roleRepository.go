package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
	"log"
	"time"
)

type RoleRepository struct {
	connection *gorm.DB
}

func NewRoleRepositoryPostgres(connection *gorm.DB) *RoleRepository {
	return &RoleRepository{connection: connection}
}

func (r *RoleRepository) Create(role models.Role) (models.Role, error) {
	result := r.connection.
		Omit("updated_at", "deleted_at").
		Create(&role)

	if err := result.Error; err != nil {
		log.Println("RoleRepository.Create(): error occured during Role creation. Passed data: ", role)
		log.Println(err)
		return models.Role{}, err
	}

	return role, nil
}

func (r *RoleRepository) GetAll() []models.Role {
	var allRoles []models.Role

	r.connection.Find(&allRoles)

	return allRoles
}

func (r *RoleRepository) GetRoleById(roleId int) (models.Role, error) {
	var foundRole models.Role

	result := r.connection.Find(&foundRole, "role_id", roleId)
	if err := result.Error; err != nil {
		log.Println("RoleRepository.GetRoleById(): error occured during Role search. Passed data: ", roleId)
		log.Println(err)
		return models.Role{}, err
	}

	return foundRole, nil
}

func (r *RoleRepository) Update(role models.Role) (models.Role, error) {
	result := r.connection.
		Omit("active", "created_at", "deleted_at").
		Model(&role).
		Updates(&role)

	if err := result.Error; err != nil {
		log.Println("RoleRepository.Update(): error occured during Role update. Passed data: ", role)
		return role, err
	}

	return role, nil
}

func (r *RoleRepository) Delete(roleId int) (bool, error) {
	roleToDelete := models.Role{
		RoleId:    roleId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := r.connection.
		Select("*").
		Omit("created_at", "updated_at", "name", "description").
		Model(&roleToDelete).
		Updates(&roleToDelete)

	if err := result.Error; err != nil {
		log.Println("RoleRepository.Delete(): error occured during Role deletion. Passed data: ", roleId)
		log.Println(err)
		return false, err
	}

	return true, nil
}
