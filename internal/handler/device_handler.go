package handler

import (
	"log"
	"net/http"

	"github.com/astanx/anime_api/internal/service"
	"github.com/google/uuid"

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
func (h *DeviceHandler) AddDeviceID(c *gin.Context) {
	deviceId := uuid.New()

	_, err := h.service.AddDeviceID(deviceId)
	if err != nil {
		log.Printf("failed to add device: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't add device"})
		return
	}

	type response struct {
		ID string `json:"id"`
	}

	c.JSON(http.StatusOK, response{ID: deviceId.String()})
}
