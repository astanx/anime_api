package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type AnimeService struct {
	repo *repository.AnimeRepo
}

func NewAnimeService(repo *repository.AnimeRepo) *AnimeService {
	return &AnimeService{repo: repo}
}

func (s *AnimeService) SearchConsumetAnime(query string) ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetAnime(query)
}

func (s *AnimeService) SearchAnilibriaAnime(query string) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaAnime(query)
}

func (s *AnimeService) SearchConsumetLatestReleases(limit int) ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetLatestReleases(limit)
}

func (s *AnimeService) SearchAnilibriaLatestReleases(limit int) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaLatestReleases(limit)
}

func (s *AnimeService) SearchAnilibriaRandomReleases(limit int) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaRandomReleases(limit)
}

func (s *AnimeService) SearchConsumetGenreReleases(genre string) ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetGenreReleases(genre)
}

func (s *AnimeService) SearchAnilibriaGenreReleases(genreID, limit int) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaGenreReleases(genreID, limit)
}

func (s *AnimeService) GetConsumetGenres() ([]model.ConsumetGenre, error) {
	return s.repo.GetConsumetGenres()
}

func (s *AnimeService) GetAnilibriaGenres() ([]model.Genre, error) {
	return s.repo.GetAnilibriaGenres()
}
