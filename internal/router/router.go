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
		// Device routes
		deviceRepository := repository.NewDeviceRepo(databases)
		deviceService := service.NewDeviceService(deviceRepository)
		deviceHandler := handler.NewDeviceHandler(deviceService)

		users := v1.Group("/users")
		{
			users.GET("/device", deviceHandler.AddDeviceID)
		}

		// Anime routes
		animeRepository := repository.NewAnimeRepo(databases)
		animeService := service.NewAnimeService(animeRepository)
		animeHandler := handler.NewAnimeHandler(animeService)

		anime := v1.Group("/anime")
		{
			// Consumet routes
			consumet := anime.Group("/consumet")
			{
				consumet.GET("/", animeHandler.SearchConsumetAnime)
				consumet.GET("/genres", animeHandler.GetConsumetGenres)
				consumet.GET("/latest", animeHandler.SearchConsumetLatestReleases)
				consumet.GET("/recommended", animeHandler.SearchConsumetRecommendedAnime)
				consumet.GET("/genre/releases", animeHandler.SearchConsumetGenreReleases)
				consumet.GET("/:id", animeHandler.GetAnimeInfoByConsumetID)
				consumet.GET("/episode/:id", animeHandler.GetConsumetEpisodeInfo)
			}

			// Anilibria routes
			anilibria := anime.Group("/anilibria")
			{
				anilibria.GET("/", animeHandler.SearchAnilibriaAnime)
				anilibria.GET("/genres", animeHandler.GetAnilibriaGenres)
				anilibria.GET("/latest", animeHandler.SearchAnilibriaLatestReleases)
				anilibria.GET("/random", animeHandler.SearchAnilibriaRandomReleases)
				anilibria.GET("/recommended", animeHandler.SearchAnilibriaRecommendedAnime)
				anilibria.GET("/genre/releases", animeHandler.SearchAnilibriaGenreReleases)
				anilibria.GET("/:id", animeHandler.GetAnimeInfoByAnilibriaID)
				anilibria.GET("/episode/:id", animeHandler.GetAnilibriaEpisodeInfo)
			}

			anime.GET("/search/:id", animeHandler.SearchAnimeByID)
		}
	}

	return r
}
