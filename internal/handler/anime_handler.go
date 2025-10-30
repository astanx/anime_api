package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type AnimeHandler struct {
	service *service.AnimeService
}

func NewAnimeHandler(s *service.AnimeService) *AnimeHandler {
	return &AnimeHandler{
		service: s,
	}
}

func (h *AnimeHandler) SearchConsumetAnime(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		log.Println("missing query param in SearchConsumetAnime")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param is required"})
		return
	}

	anime, err := h.service.SearchConsumetAnime(query)
	if err != nil {
		log.Printf("SearchConsumetAnime: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": anime})
}

func (h *AnimeHandler) SearchAnilibriaAnime(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		log.Println("missing query param in SearchAnilibriaAnime")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param is required"})
		return
	}

	anime, err := h.service.SearchAnilibriaAnime(query)
	if err != nil {
		log.Printf("SearchAnilibriaAnime: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": anime})
}

func (h *AnimeHandler) GetAnilibriaGenres(c *gin.Context) {
	genres, err := h.service.GetAnilibriaGenres()
	if err != nil {
		log.Printf("GetAnilibriaGenres: failed to get genres: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get genres"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) GetConsumetGenres(c *gin.Context) {
	genres, err := h.service.GetConsumetGenres()
	if err != nil {
		log.Printf("GetConsumetGenres: failed to get genres: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get genres"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchAnilibriaGenreReleases(c *gin.Context) {
	genreIdStr := c.Query("genre")
	if genreIdStr == "" {
		log.Println("missing genre param in SearchAnilibriaGenreReleases")
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre param is required"})
		return
	}
	genreId, err := strconv.Atoi(genreIdStr)
	if err != nil {
		log.Printf("SearchAnilibriaGenreReleases: invalid genre: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre must be an integer"})
		return
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		log.Println("missing limit param in SearchAnilibriaGenreReleases")
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit param is required"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("SearchAnilibriaGenreReleases: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}
	genres, err := h.service.SearchAnilibriaGenreReleases(genreId, limit)
	if err != nil {
		log.Printf("SearchAnilibriaGenreReleases: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchConsumetGenreReleases(c *gin.Context) {
	genre := c.Query("genre")
	if genre == "" {
		log.Println("missing genre param in SearchConsumetGenreReleases")
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre param is required"})
		return
	}
	genres, err := h.service.SearchConsumetGenreReleases(genre)
	if err != nil {
		log.Printf("SearchConsumetGenreReleases: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchAnilibriaLatestReleases(c *gin.Context) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		log.Println("missing limit param in SearchAnilibriaLatestReleases")
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit param is required"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("SearchAnilibriaLatestReleases: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}
	genres, err := h.service.SearchAnilibriaLatestReleases(limit)
	if err != nil {
		log.Printf("SearchAnilibriaLatestReleases: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchConsumetLatestReleases(c *gin.Context) {
	genres, err := h.service.SearchConsumetLatestReleases()
	if err != nil {
		log.Printf("SearchConsumetLatestReleases: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchAnilibriaRandomReleases(c *gin.Context) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		log.Println("missing limit param in SearchAnilibriaRandomReleases")
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit param is required"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("SearchAnilibriaRandomReleases: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}
	genres, err := h.service.SearchAnilibriaRandomReleases(limit)
	if err != nil {
		log.Printf("SearchAnilibriaRandomReleases: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": genres})
}

func (h *AnimeHandler) SearchAnimeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in GetSearchAnimeByID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	anime, err := h.service.SearchAnimeByID(id)
	if err != nil {
		log.Printf("GetSearchAnimeByID: failed to get anime: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get anime"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": anime})
}

func (h *AnimeHandler) GetAnimeInfoByConsumetID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in GetAnimeInfoByConsumetID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	anime, err := h.service.GetAnimeInfoByConsumetID(id)
	if err != nil {
		log.Printf("GetAnimeInfoByConsumetID: failed to get anime info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get anime info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": anime})
}

func (h *AnimeHandler) GetAnimeInfoByAnilibriaID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in GetAnimeInfoByAnilibriaID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	anime, err := h.service.GetAnimeInfoByAnilibriaID(id)
	if err != nil {
		log.Printf("GetAnimeInfoByAnilibriaID: failed to get anime info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get anime info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": anime})
}

func (h *AnimeHandler) GetAnilibriaEpisodeInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in GetAnilibriaEpisodeInfo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	episode, err := h.service.GetAnilibriaEpisodeInfo(id)
	if err != nil {
		log.Printf("GetAnilibriaEpisodeInfo: failed to get episode info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get episode info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": episode})
}

func (h *AnimeHandler) GetConsumetEpisodeInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in GetConsumetEpisodeInfo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}
	title := c.Query("title")
	ordinalStr := c.Query("ordinal")
	dub := c.DefaultQuery("dub", "false")

	var ordinal int
	if ordinalStr == "" {
		ordinal = -1
	} else {
		var err error
		ordinal, err = strconv.Atoi(ordinalStr)
		if err != nil {
			log.Printf("SearchAnilibriaRecommendedAnime: invalid limit: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
			return
		}
	}

	episode, err := h.service.GetConsumetEpisodeInfo(id, title, ordinal, dub)
	if err != nil {
		log.Printf("GetConsumetEpisodeInfo: failed to get episode info: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get episode info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": episode})
}

func (h *AnimeHandler) SearchConsumetRecommendedAnime(c *gin.Context) {
	anime, err := h.service.SearchConsumetRecommendedAnime()
	if err != nil {
		log.Printf("SearchConsumetRecommendedAnime: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": anime})
}

func (h *AnimeHandler) SearchAnilibriaRecommendedAnime(c *gin.Context) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "14"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("SearchAnilibriaRecommendedAnime: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}

	anime, err := h.service.SearchAnilibriaRecommendedAnime(limit)
	if err != nil {
		log.Printf("SearchAnilibriaRecommendedAnime: failed to get releases: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": anime})
}
