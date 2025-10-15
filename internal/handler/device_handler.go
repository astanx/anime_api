package handler

import (
	"log"
	"net/http"

	"github.com/astanx/anime_api/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *service.DeviceService
}

func NewDeviceHandler(s *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		service: s,
	}
}

func (h *DeviceHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUserByDeviceID("aaa")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *DeviceHandler) AddDeviceID(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	_, err := h.service.AddDeviceID(req.DeviceID)
	if err != nil {
		log.Println("failed to add device: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant add device"})
		return
	}
	c.Status(http.StatusOK)
}
