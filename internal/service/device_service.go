package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
	"github.com/google/uuid"
)

type DeviceService struct {
	repo *repository.DeviceRepo
}

func NewDeviceService(repo *repository.DeviceRepo) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) AddDeviceID(deviceID uuid.UUID) (model.User, error) {
	return s.repo.AddDeviceID(deviceID)
}
