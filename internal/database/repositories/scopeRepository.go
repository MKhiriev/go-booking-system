package repositories

import (
	"go-booking-system/internal/models"
	"gorm.io/gorm"
	"log"
	"time"
)

type ScopeRepository struct {
	connection *gorm.DB
}

func NewScopeRepositoryPostgres(connection *gorm.DB) *ScopeRepository {
	return &ScopeRepository{connection: connection}
}

func (r *ScopeRepository) Create(scope models.Scope) (models.Scope, error) {
	result := r.connection.
		Omit("updated_at", "deleted_at").
		Create(&scope)

	if err := result.Error; err != nil {
		log.Println("ScopeRepository.Create(): error occured during Scope creation. Passed data: ", scope)
		log.Println(err)
		return models.Scope{}, err
	}

	return scope, nil
}

func (r *ScopeRepository) GetAll() []models.Scope {
	var allScopes []models.Scope

	r.connection.Find(&allScopes)

	return allScopes
}

func (r *ScopeRepository) GetScopeById(scopeId int) (models.Scope, error) {
	var foundScope models.Scope

	result := r.connection.Find(&foundScope, "scope_id", scopeId)
	if err := result.Error; err != nil {
		log.Println("ScopeRepository.GetRoomById(): error occured during Scope search. Passed data: ", scopeId)
		log.Println(err)
		return models.Scope{}, err
	}

	return foundScope, nil
}

func (r *ScopeRepository) Update(scope models.Scope) (models.Scope, error) {
	result := r.connection.
		Omit("active", "created_at", "deleted_at").
		Model(&scope).
		Updates(&scope)

	if err := result.Error; err != nil {
		log.Println("ScopeRepository.Update(): error occured during Scope update. Passed data: ", scope)
		return scope, err
	}

	return scope, nil
}

func (r *ScopeRepository) Delete(roleId int) (bool, error) {
	scopeToDelete := models.Scope{
		ScopeId:   roleId,
		Active:    false,
		DeletedAt: time.Now(),
	}

	result := r.connection.
		Select("*").
		Omit("created_at", "updated_at", "name", "description").
		Model(&scopeToDelete).
		Updates(&scopeToDelete)

	if err := result.Error; err != nil {
		log.Println("ScopeRepository.Delete(): error occured during Scope deletion. Passed data: ", roleId)
		log.Println(err)
		return false, err
	}

	return true, nil
}
