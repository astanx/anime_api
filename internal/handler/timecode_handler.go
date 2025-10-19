package handler

import (
	"log"
	"net/http"

	"github.com/astanx/anime_api/internal/model"
	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type TimecodeHandler struct {
	service *service.TimecodeService
}

func NewTimecodeHandler(s *service.TimecodeService) *TimecodeHandler {
	return &TimecodeHandler{
		service: s,
	}
}

func (h *TimecodeHandler) GetAllTimecodes(c *gin.Context) {
	deviceID := c.Query("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	timecodes, err := h.service.GetAllTimecodes(deviceID)
	if err != nil {
		log.Printf("failed to get all timecodes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get timecodes"})
		return
	}

	c.JSON(http.StatusOK, timecodes)
}

func (h *TimecodeHandler) GetTimecode(c *gin.Context) {
	deviceID := c.Query("deviceID")
	episodeID := c.Query("episodeID")

	if deviceID == "" || episodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID and episodeID are required"})
		return
	}

	timecode, err := h.service.GetTimecode(deviceID, episodeID)
	if err != nil {
		log.Printf("failed to get timecode: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get timecode"})
		return
	}

	if timecode == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "timecode not found"})
		return
	}

	c.JSON(http.StatusOK, timecode)
}

func (h *TimecodeHandler) AddOrUpdateTimecode(c *gin.Context) {
	deviceID := c.Query("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	var timecode model.Timecode
	if err := c.ShouldBindJSON(&timecode); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.AddOrUpdateTimecode(deviceID, timecode); err != nil {
		log.Printf("failed to add/update timecode: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add/update timecode"})
		return
	}

	c.Status(http.StatusNoContent)
}
