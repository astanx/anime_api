package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type HistoryService struct {
	repo *repository.HistoryRepo
}

func NewHistoryService(repo *repository.HistoryRepo) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) AddHistory(deviceID string, history model.History) error {
	return s.repo.AddHistory(deviceID, history)
}

func (s *HistoryService) GetAllHistory(deviceID string) ([]model.History, error) {
	return s.repo.GetAllHistory(deviceID)
}

func (s *HistoryService) GetHistory(deviceID string, page, limit int) (model.PaginatedHistory, error) {
	return s.repo.GetHistory(deviceID, page, limit)
}
