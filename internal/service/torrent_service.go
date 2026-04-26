package service

import (
	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/repository"
)

type TorrentService struct {
	repo *repository.TorrentRepo
}

func NewTorrentService(repo *repository.TorrentRepo) *TorrentService {
	return &TorrentService{repo: repo}
}

func (s *TorrentService) SearchMALAnime(query string, page int) (model.PaginatedSearchAnime, error) {
	return s.repo.SearchMALAnime(query, page)
}

func (s *TorrentService) SearchMALRecommendedAnime(limit, page int) ([]model.SearchAnime, error) {
	return s.repo.SearchMALRecommendedAnime(limit, page)
}

func (s *TorrentService) SearchMALLatestReleases(limit, page int) ([]model.SearchAnime, error) {
	return s.repo.SearchMALLatestReleases(limit, page)
}

func (s *TorrentService) SearchMALById(id string) (model.Anime, error) {
	return s.repo.SearchMALById(id)
}

func (s *TorrentService) SearchMALByEpisodeId(animeId, episodeId string) (model.Episode, error) {
	return s.repo.SearchMALByEpisodeId(animeId, episodeId)
}
