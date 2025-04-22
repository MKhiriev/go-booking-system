package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
)

type PermissionService struct {
	repository database.PermissionRepository
}

func NewPermissionService(repository database.PermissionRepository) *PermissionService {
	return &PermissionService{repository: repository}
}

func (r *PermissionService) Create(permission models.Permission) (models.Permission, error) {
	return r.repository.Create(permission)
}

func (r *PermissionService) GetAll() []models.Permission {
	return r.repository.GetAll()
}

func (r *PermissionService) GetPermissionsByRoleId(roleId int) ([]models.Permission, error) {
	return r.repository.GetPermissionsByRoleId(roleId)
}

func (r *PermissionService) GetPermissionsByRouteId(routeId int) ([]models.Permission, error) {
	return r.repository.GetPermissionsByRouteId(routeId)
}

func (r *PermissionService) GetPermissionsByRoleIdAndRouteId(roleId int, routeId int) ([]models.Permission, error) {
	return r.repository.GetPermissionsByRoleIdAndRouteId(roleId, routeId)
}

func (r *PermissionService) Update(permission models.Permission) (models.Permission, error) {
	return r.repository.Update(permission)
}

func (r *PermissionService) Delete(roleId int, routeId int) (bool, error) {
	return r.repository.Delete(roleId, routeId)
}
