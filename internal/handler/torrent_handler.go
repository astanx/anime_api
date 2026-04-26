package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type TorrentHandler struct {
	service *service.TorrentService
}

func NewTorrentHandler(s *service.TorrentService) *TorrentHandler {
	return &TorrentHandler{
		service: s,
	}
}

func (h *TorrentHandler) SearchMALAnime(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		log.Println("missing query param in SearchMALAnime")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param is required"})
		return
	}
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("SearchMALAnime: invalid page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be an integer"})
		return
	}

	anime, err := h.service.SearchMALAnime(query, pageInt)
	if err != nil {
		log.Printf("SearchMALAnime: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": anime})
}

func (h *TorrentHandler) SearchMALRecommendedAnime(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("SearchMALRecommendedAnime: invalid page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be an integer"})
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("SearchMALRecommendedAnime: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}

	data, err := h.service.SearchMALRecommendedAnime(limitInt, pageInt)
	if err != nil {
		log.Printf("SearchMALRecommendedAnime: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": data})
}

func (h *TorrentHandler) SearchMALLatestReleases(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("SearchMALLatestReleases: invalid page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be an integer"})
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("SearchMALLatestReleases: invalid limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}

	data, err := h.service.SearchMALLatestReleases(pageInt, limitInt)
	if err != nil {
		log.Printf("SearchMALLatestReleases: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": data})
}

func (h *TorrentHandler) SearchMALById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Println("missing id param in SearchMALById")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	data, err := h.service.SearchMALById(id)
	if err != nil {
		log.Printf("SearchMALById: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": data})
}

func (h *TorrentHandler) SearchMALByEpisodeId(c *gin.Context) {
	animeId := c.Param("id")
	if animeId == "" {
		log.Println("missing id param in SearchMALByEpisodeId")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}
	episodeId := c.Param("episodeId")
	if episodeId == "" {
		log.Println("missing episodeId param in SearchMALByEpisodeId")
		c.JSON(http.StatusBadRequest, gin.H{"error": "episodeId param is required"})
		return
	}

	data, err := h.service.SearchMALByEpisodeId(animeId, episodeId)
	if err != nil {
		log.Printf("SearchMALByEpisodeId: failed to search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to search"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": data})
}
