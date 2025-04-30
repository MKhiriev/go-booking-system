package services

import (
	"go-booking-system/internal/database"
	"go-booking-system/internal/models"
)

type RoleService struct {
	repository database.RoleRepository
}

func NewRoleService(repository database.RoleRepository) *RoleService {
	return &RoleService{repository: repository}
}

func (r *RoleService) Create(role models.Role) (models.Role, error) {
	return r.repository.Create(role)
}

func (r *RoleService) GetAll() []models.Role {
	return r.repository.GetAll()
}

func (r *RoleService) GetRoleById(roleId int) (models.Role, error) {
	return r.repository.GetRoleById(roleId)
}

func (r *RoleService) Update(role models.Role) (models.Role, error) {
	return r.repository.Update(role)
}

func (r *RoleService) Delete(roleId int) (bool, error) {
	return r.repository.Delete(roleId)
}
