package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type TimecodeService struct {
	repo *repository.TimecodeRepo
}

func NewTimecodeService(repo *repository.TimecodeRepo) *TimecodeService {
	return &TimecodeService{repo: repo}
}

func (s *TimecodeService) GetAllTimecodes(deviceID string) ([]model.Timecode, error) {
	return s.repo.GetAllTimecodes(deviceID)
}

func (s *TimecodeService) GetTimecode(deviceID, episodeID string) (*model.Timecode, error) {
	return s.repo.GetTimecode(deviceID, episodeID)
}

func (s *TimecodeService) AddOrUpdateTimecode(deviceID string, timecode model.Timecode) error {
	return s.repo.AddTimecode(deviceID, timecode)
}
