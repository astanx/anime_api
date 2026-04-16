package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/astanx/anime_api/internal/service"
	"github.com/gin-gonic/gin"
)

type MALHandler struct {
	service *service.MALService
}

func NewMALHandler(s *service.MALService) *MALHandler {
	return &MALHandler{
		service: s,
	}
}

func (h *MALHandler) ExportMALList(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	malList, err := h.service.ExportMALList(deviceID)
	if err != nil {
		log.Printf("failed to export MAL list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't export MAL list"})
		return
	}

	c.JSON(http.StatusOK, malList)
}

func (h *MALHandler) ImportMALList(c *gin.Context) {
	deviceID := c.GetString("deviceID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceID is required"})
		return
	}

	var malList struct {
		MalList string `json:"malList"`
	}

	if err := c.ShouldBindJSON(&malList); err != nil {
		log.Printf("invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	count, err := h.service.ImportMALList(deviceID, malList.MalList)
	if err != nil {
		log.Printf("failed to import MAL list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't import MAL list", "message": "Error occured during MAL import."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MAL list imported successfully with " + fmt.Sprintf("%d", count) + " titles."})
}
