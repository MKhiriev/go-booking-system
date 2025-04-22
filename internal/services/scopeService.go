package services

import (
	"humoBooking/internal/database"
	"humoBooking/internal/models"
)

type ScopeService struct {
	repository database.ScopeRepository
}

func NewScopeService(repository database.ScopeRepository) *ScopeService {
	return &ScopeService{repository: repository}
}

func (r *ScopeService) Create(scope models.Scope) (models.Scope, error) {
	return r.repository.Create(scope)
}

func (r *ScopeService) GetAll() []models.Scope {
	return r.repository.GetAll()
}

func (r *ScopeService) GetScopeById(scopeId int) (models.Scope, error) {
	return r.repository.GetScopeById(scopeId)
}

func (r *ScopeService) Update(scope models.Scope) (models.Scope, error) {
	return r.repository.Update(scope)
}

func (r *ScopeService) Delete(scopeId int) (bool, error) {
	return r.repository.Delete(scopeId)
}
