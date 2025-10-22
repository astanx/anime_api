package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type CollectionService struct {
	repo *repository.CollectionRepo
}

func NewCollectionService(repo *repository.CollectionRepo) *CollectionService {
	return &CollectionService{repo: repo}
}

func (s *CollectionService) AddCollection(deviceID string, collection model.Collection) error {
	return s.repo.AddCollection(deviceID, collection)
}

func (s *CollectionService) RemoveCollection(deviceID, animeID, collectionType string) error {
	return s.repo.RemoveCollection(deviceID, animeID, collectionType)
}

func (s *CollectionService) GetAllCollections(deviceID string) ([]model.Collection, error) {
	return s.repo.GetAllCollections(deviceID)
}

func (s *CollectionService) GetCollections(deviceID, T string, page, limit int) (model.PaginatedCollections, error) {
	return s.repo.GetCollections(deviceID, T, page, limit)
}
