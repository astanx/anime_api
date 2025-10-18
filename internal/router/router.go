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
			users.GET("/device", deviceHandler.AddDeviceID)
		}
	}

	{
		animeRepository := repository.NewAnimeRepo(databases)
		animeService := service.NewAnimeService(animeRepository)
		animeHandler := handler.NewAnimeHandler(animeService)

		anime := v1.Group("/anime")
		{
			consumet := anime.Group("/consumet")
			{
				consumet.GET("/", animeHandler.SearchConsumetAnime)
				consumet.GET("/genres", animeHandler.GetConsumetGenres)
				consumet.GET("/latest", animeHandler.SearchConsumetLatestReleases)
				consumet.GET("/genre/releases", animeHandler.SearchConsumetGenreReleases)
			}
			anilibria := anime.Group("/anilibria")
			{
				anilibria.GET("/", animeHandler.SearchAnilibriaAnime)
				anilibria.GET("/genres", animeHandler.GetAnilibriaGenres)
				anilibria.GET("/latest", animeHandler.SearchAnilibriaLatestReleases)
				anilibria.GET("/random", animeHandler.SearchAnilibriaRandomReleases)
				anilibria.GET("/genre/releases", animeHandler.SearchAnilibriaGenreReleases)
			}
		}
	}

	return r
}
