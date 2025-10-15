package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type DeviceService struct {
	repo *repository.DeviceRepo
}

func NewDeviceService(repo *repository.DeviceRepo) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) GetUserByDeviceID(deviceID string) (model.User, error) {
	return s.repo.GetByDeviceID(deviceID)
}

func (s *DeviceService) AddDeviceID(deviceID string) (model.User, error) {
	return s.repo.AddDeviceID(deviceID)
}
