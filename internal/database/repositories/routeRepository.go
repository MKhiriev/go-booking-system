package repositories

import (
	"go-booking-system/internal/models"
	"gorm.io/gorm"
	"log"
	"time"
)

type RouteRepository struct {
	connection *gorm.DB
}

func NewRouteRepositoryPostgres(connection *gorm.DB) *RouteRepository {
	return &RouteRepository{connection: connection}
}

func (r *RouteRepository) Create(route models.Route) (models.Route, error) {
	result := r.connection.
		Omit("updated_at", "deleted_at").
		Create(&route)

	if err := result.Error; err != nil {
		log.Println("RouteRepository.Create(): error occured during Route creation. Passed data: ", route)
		log.Println(err)
		return models.Route{}, err
	}

	return route, nil
}

func (r *RouteRepository) GetAll() []models.Route {
	var allRoutes []models.Route

	r.connection.Find(&allRoutes)

	return allRoutes
}

func (r *RouteRepository) GetRouteById(routeId int) (models.Route, error) {
	var foundRoute models.Route

	result := r.connection.Find(&foundRoute, "route_id", routeId)
	if err := result.Error; err != nil {
		log.Println("RouteRepository.GetRouteById(): error occured during Route search. Passed data: ", routeId)
		log.Println(err)
		return models.Route{}, err
	}

	return foundRoute, nil
}

func (r *RouteRepository) GetRouteByURL(url string) (models.Route, error) {
	var foundRoute models.Route

	result := r.connection.Find(&foundRoute, "url", url)
	if err := result.Error; err != nil {
		log.Println("RouteRepository.GetRouteByURL(): error occured during Route search. Passed data: ", url)
		log.Println(err)
		return models.Route{}, err
	}

	return foundRoute, nil
}

func (r *RouteRepository) Update(route models.Route) (models.Route, error) {
	result := r.connection.
		Omit("active", "created_at", "deleted_at").
		Model(&route).
		Updates(&route)

	if err := result.Error; err != nil {
		log.Println("RouteRepository.Update(): error occured during Route update. Passed data: ", route)
		return route, err
	}

	return route, nil
}

func (r *RouteRepository) Delete(routeId int) (bool, error) {
	routeToDelete := models.Route{
		RouteId:   routeId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := r.connection.
		Select("*").
		Omit("created_at", "updated_at", "url", "description").
		Model(&routeToDelete).
		Updates(&routeToDelete)

	if err := result.Error; err != nil {
		log.Println("RouteRepository.Delete(): error occured during Route deletion. Passed data: ", routeId)
		log.Println(err)
		return false, err
	}

	return true, nil
}
