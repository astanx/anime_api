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

// Search
func (s *AnimeService) SearchConsumetAnime(query string) ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetAnime(query)
}

func (s *AnimeService) SearchAnilibriaAnime(query string) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaAnime(query)
}

func (s *AnimeService) SearchConsumetRecommendedAnime() ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetRecommendedAnime()
}

func (s *AnimeService) SearchAnilibriaRecommendedAnime(limit int) ([]model.SearchAnime, error) {
	return s.repo.SearchAnilibriaRecommendedAnime(limit)
}

func (s *AnimeService) SearchConsumetLatestReleases() ([]model.SearchAnime, error) {
	return s.repo.SearchConsumetLatestReleases()
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

// Get genres
func (s *AnimeService) GetConsumetGenres() ([]string, error) {
	return s.repo.GetConsumetGenres()
}

func (s *AnimeService) GetAnilibriaGenres() ([]model.Genre, error) {
	return s.repo.GetAnilibriaGenres()
}

// Get anime info
func (s *AnimeService) SearchAnimeByID(id string) (model.SearchAnime, error) {
	return s.repo.SearchAnimeByID(id)
}

func (s *AnimeService) GetAnimeInfoByConsumetID(id string) (model.Anime, error) {
	return s.repo.GetAnimeInfoByConsumetID(id)
}

func (s *AnimeService) GetAnimeInfoByAnilibriaID(id string) (model.Anime, error) {
	return s.repo.GetAnimeInfoByAnilibriaID(id)
}

// Get episode info
func (s *AnimeService) GetAnilibriaEpisodeInfo(id string) (model.Episode, error) {
	return s.repo.GetAnilibriaEpisodeInfo(id)
}

func (s *AnimeService) GetConsumetEpisodeInfo(id string, title string, ordinal int, dub string) (model.Episode, error) {
	return s.repo.GetConsumetEpisodeInfo(id, title, ordinal, dub)
}
