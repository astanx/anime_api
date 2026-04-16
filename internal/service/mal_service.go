package service

import (
	"github.com/astanx/anime_api/internal/repository"
)

type MALService struct {
	repo *repository.MALRepo
}

func NewMALService(repo *repository.MALRepo) *MALService {
	return &MALService{repo: repo}
}

func (s *MALService) ExportMALList(deviceID string) (string, error) {
	return s.repo.ExportMALList(deviceID)
}

func (s *MALService) ImportMALList(deviceID, malList string) (int, error) {
	return s.repo.ImportMALList(deviceID, malList)
}
