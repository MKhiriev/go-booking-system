package database

import (
	"go-booking-system/internal/database/repositories"
	"go-booking-system/internal/models"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	BookingRepository
	RoomRepository
	UserRepository
	RoleRepository
	RouteRepository
	ScopeRepository
	PermissionRepository
}

func NewDatabase(conn *gorm.DB) *Database {
	return &Database{
		BookingRepository:    repositories.NewBookingRepositoryPostgres(conn),
		RoomRepository:       repositories.NewRoomRepositoryPostgres(conn),
		UserRepository:       repositories.NewUserRepositoryPostgres(conn),
		RoleRepository:       repositories.NewRoleRepositoryPostgres(conn),
		RouteRepository:      repositories.NewRouteRepositoryPostgres(conn),
		ScopeRepository:      repositories.NewScopeRepositoryPostgres(conn),
		PermissionRepository: repositories.NewPermissionRepositoryPostgres(conn),
	}
}

type BookingRepository interface {
	Create(booking models.Booking) (models.Booking, error)
	GetAll() []models.Booking
	GetBookingById(bookingId int) (models.Booking, error)
	GetBookingsByRoomId(roomId int) ([]models.Booking, error)
	GetBookingsByRoomIdAndBookingTime(roomId int, dateTimeStart time.Time, dateTimeEnd time.Time) ([]models.Booking, error)
	Update(booking models.Booking) (models.Booking, error)
	Delete(bookingId int) (bool, error)
}

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetAll() []models.User
	GetUserById(userId int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(userId int) (bool, error)
	UpdatePassword(user models.User) (models.User, error)
	UpdateUsername(user models.User) (models.User, error)
	UpdateUserRole(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type RoomRepository interface {
	Create(room models.Room) (models.Room, error)
	GetAll() []models.Room
	GetRoomById(roomId int) (models.Room, error)
	Update(room models.Room) (models.Room, error)
	Delete(roomId int) (bool, error)
}

type RoleRepository interface {
	Create(role models.Role) (models.Role, error)
	GetAll() []models.Role
	GetRoleById(roleId int) (models.Role, error)
	Update(role models.Role) (models.Role, error)
	Delete(roleId int) (bool, error)
}

type RouteRepository interface {
	Create(route models.Route) (models.Route, error)
	GetAll() []models.Route
	GetRouteById(routeId int) (models.Route, error)
	GetRouteByURL(url string) (models.Route, error)
	Update(route models.Route) (models.Route, error)
	Delete(routeId int) (bool, error)
}

type ScopeRepository interface {
	Create(scope models.Scope) (models.Scope, error)
	GetAll() []models.Scope
	GetScopeById(scopeId int) (models.Scope, error)
	Update(scope models.Scope) (models.Scope, error)
	Delete(scopeId int) (bool, error)
}

type PermissionRepository interface {
	Create(permission models.Permission) (models.Permission, error)
	GetAll() []models.Permission
	GetPermissionsByRoleId(roleId int) ([]models.Permission, error)
	GetPermissionsByRouteId(routeId int) ([]models.Permission, error)
	GetPermissionsByRoleIdAndRouteId(roleId int, routeId int) ([]models.Permission, error)
	Update(permission models.Permission) (models.Permission, error)
	Delete(roleId int, routeId int) (bool, error)
}
