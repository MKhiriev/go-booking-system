package repositories

import (
	"gorm.io/gorm"
	"humoBooking/internal/models"
	"log"
	"time"
)

type PermissionRepository struct {
	connection *gorm.DB
}

func NewPermissionRepositoryPostgres(connection *gorm.DB) *PermissionRepository {
	return &PermissionRepository{connection: connection}
}

func (r *PermissionRepository) Create(permission models.Permission) (models.Permission, error) {
	result := r.connection.
		Omit("updated_at", "deleted_at").
		Create(&permission)

	if err := result.Error; err != nil {
		log.Println("PermissionRepository.Create(): error occured during Permission creation. Passed data: ", permission)
		log.Println(err)
		return models.Permission{}, err
	}

	return permission, nil
}

func (r *PermissionRepository) GetAll() []models.Permission {
	var allPermissions []models.Permission

	r.connection.Find(&allPermissions)

	return allPermissions
}

func (r *PermissionRepository) GetPermissionsByRoleId(roleId int) ([]models.Permission, error) {
	var foundPermissions []models.Permission

	result := r.connection.Find(&foundPermissions, "role_id", roleId)
	if err := result.Error; err != nil {
		log.Println("PermissionRepository.GetPermissionByRoleId(): error occured during Permissions search. Passed data: ", roleId)
		log.Println(err)
		return []models.Permission{}, err
	}

	return foundPermissions, nil
}

func (r *PermissionRepository) GetPermissionsByRouteId(routeId int) ([]models.Permission, error) {
	var foundPermissions []models.Permission

	result := r.connection.Find(&foundPermissions, "route_id", routeId)
	if err := result.Error; err != nil {
		log.Println("PermissionRepository.GetPermissionsByRouteId(): error occured during Permissions search. Passed data: ", routeId)
		log.Println(err)
		return []models.Permission{}, err
	}

	return foundPermissions, nil
}

func (r *PermissionRepository) Update(permission models.Permission) (models.Permission, error) {
	result := r.connection.
		Omit("active", "created_at", "deleted_at").
		Model(&permission).
		Updates(&permission)

	if err := result.Error; err != nil {
		log.Println("PermissionRepository.Update(): error occured during Permission update. Passed data: ", permission)
		return permission, err
	}

	return permission, nil
}

func (r *PermissionRepository) Delete(roleId int, routeId int) (bool, error) {
	permissionToDelete := models.Permission{
		RoleId:    roleId,
		RouteId:   routeId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := r.connection.
		Select("*").
		Omit("created_at", "updated_at", "name", "description").
		Model(&permissionToDelete).
		Updates(&permissionToDelete)

	if err := result.Error; err != nil {
		log.Printf("PermissionRepository.Delete(): error occured during Permission deletion. Passed data: RoleId=%d RouteId=%d", roleId, roleId)
		log.Println(err)
		return false, err
	}

	return true, nil
}
