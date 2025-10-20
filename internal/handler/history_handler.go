package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service *service.HistoryService
}

func NewHistoryHandler(s *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		service: s,
	}
}

func (h *HistoryHandler) AddHistory(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	var history model.History
	if err := c.ShouldBindJSON(&history); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.AddHistory(deviceID, history); err != nil {
		log.Printf("failed to add history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add history"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *HistoryHandler) GetAllHistory(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	history, err := h.service.GetAllHistory(deviceID)
	if err != nil {
		log.Printf("failed to get all history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
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

	history, err := h.service.GetHistory(deviceID, page, limit)
	if err != nil {
		log.Printf("failed to get paginated history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
