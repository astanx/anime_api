// internal/router/router.go
package router

import (
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/handler"
	"github.com/astanx/anime_api/internal/middleware"
	"github.com/astanx/anime_api/internal/repository"
	"github.com/astanx/anime_api/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(databases *db.DB) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logging())

	v1 := r.Group("/api/v1")
	{

		deviceRepository := repository.NewDeviceRepo(databases)
		deviceService := service.NewDeviceService(deviceRepository)
		deviceHandler := handler.NewDeviceHandler(deviceService)

		users := v1.Group("/users")
		{
			users.POST("/device", deviceHandler.AddDeviceID)
		}
	}

	return r
}
