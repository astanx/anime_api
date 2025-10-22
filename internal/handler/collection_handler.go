package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	service *service.CollectionService
}

func NewCollectionHandler(s *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		service: s,
	}
}

func (h *CollectionHandler) AddCollection(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	var collection model.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.AddCollection(deviceID, collection); err != nil {
		log.Printf("failed to add collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add collection"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CollectionHandler) RemoveCollection(c *gin.Context) {
	var req struct {
		AnimeID        string `json:"anime_id"`
		CollectionType string `json:"type"`
	}

	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.AnimeID == "" || req.CollectionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "anime_id, and type are required"})
		return
	}

	if err := h.service.RemoveCollection(deviceID, req.AnimeID, req.CollectionType); err != nil {
		log.Printf("failed to remove collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't remove collection"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CollectionHandler) GetAllCollections(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	collections, err := h.service.GetAllCollections(deviceID)
	if err != nil {
		log.Printf("failed to get all collections: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get collections"})
		return
	}

	c.JSON(http.StatusOK, collections)
}

func (h *CollectionHandler) GetCollections(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	T := c.Query("type")
	if T == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type query is required"})
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

	collections, err := h.service.GetCollections(deviceID, T, page, limit)
	if err != nil {
		log.Printf("failed to get paginated collections: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get collections"})
		return
	}

	c.JSON(http.StatusOK, collections)
}
