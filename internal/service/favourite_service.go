package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type FavouriteService struct {
	repo *repository.FavouriteRepo
}

func NewFavouriteService(repo *repository.FavouriteRepo) *FavouriteService {
	return &FavouriteService{repo: repo}
}

func (s *FavouriteService) AddFavourite(deviceID string, favourite model.Favourite) error {
	return s.repo.AddFavourite(deviceID, favourite)
}

func (s *FavouriteService) RemoveFavourite(deviceID string, favourite model.Favourite) error {
	return s.repo.RemoveFavourite(deviceID, favourite)
}

func (s *FavouriteService) GetAllFavourites(deviceID string) ([]model.Favourite, error) {
	return s.repo.GetAllFavourites(deviceID)
}

func (s *FavouriteService) GetFavourites(deviceID string, page, limit int) (model.PaginatedFavourites, error) {
	return s.repo.GetFavourites(deviceID, page, limit)
}

func (s *FavouriteService) GetFavouriteForAnime(deviceID, animeID string) (model.Favourite, error) {
	return s.repo.GetFavouriteForAnime(deviceID, animeID)
}
