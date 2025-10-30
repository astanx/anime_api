package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type FavouriteHandler struct {
	service *service.FavouriteService
}

func NewFavouriteHandler(s *service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{
		service: s,
	}
}

func (h *FavouriteHandler) AddFavourite(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	var favourite model.Favourite
	if err := c.ShouldBindJSON(&favourite); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.AddFavourite(deviceID, favourite); err != nil {
		log.Printf("failed to add favourite: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add favourite"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *FavouriteHandler) RemoveFavourite(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}
	var favourite model.Favourite
	if err := c.ShouldBindJSON(&favourite); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.RemoveFavourite(deviceID, favourite); err != nil {
		log.Printf("failed to remove favourite: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't remove favourite"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *FavouriteHandler) GetAllFavourites(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	favourites, err := h.service.GetAllFavourites(deviceID)
	if err != nil {
		log.Printf("failed to get all favourites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get favourites"})
		return
	}

	c.JSON(http.StatusOK, favourites)
}

func (h *FavouriteHandler) GetFavourites(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	favourites, err := h.service.GetFavourites(deviceID, page, limit)
	if err != nil {
		log.Printf("failed to get paginated favourites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get favourites"})
		return
	}

	c.JSON(http.StatusOK, favourites)
}

func (h *FavouriteHandler) GetFavouriteForAnime(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	animeID := c.Query("animeID")
	if animeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "animeID is required"})
		return
	}

	favourite, err := h.service.GetFavouriteForAnime(deviceID, animeID)
	if err != nil {
		log.Printf("failed to get favourite for anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get favourite"})
		return
	}

	c.JSON(http.StatusOK, favourite)
}
