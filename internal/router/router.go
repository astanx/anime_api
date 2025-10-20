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
	r.Use(middleware.DeviceMiddleware())

	v1 := r.Group("/api/v1")
	{
		// Device routes
		deviceRepo := repository.NewDeviceRepo(databases)
		deviceService := service.NewDeviceService(deviceRepo)
		deviceHandler := handler.NewDeviceHandler(deviceService)

		users := v1.Group("/users")
		{
			users.GET("/device", deviceHandler.AddDeviceID)
		}

		// Timecode routes
		timecodeRepo := repository.NewTimecodeRepo(databases)
		timecodeService := service.NewTimecodeService(timecodeRepo)
		timecodeHandler := handler.NewTimecodeHandler(timecodeService)

		timecodes := v1.Group("/timecode")
		{
			timecodes.GET("", timecodeHandler.GetTimecode)
			timecodes.GET("/all", timecodeHandler.GetAllTimecodes)
			timecodes.POST("", timecodeHandler.AddOrUpdateTimecode)
		}

		// History routes
		historyRepo := repository.NewHistoryRepo(databases)
		historyService := service.NewHistoryService(historyRepo)
		historyHandler := handler.NewHistoryHandler(historyService)

		history := v1.Group("/history")
		{
			history.POST("", historyHandler.AddHistory)
			history.GET("", historyHandler.GetHistory)
			history.GET("/all", historyHandler.GetAllHistory)
		}

		// Favourite routes
		favRepo := repository.NewFavouriteRepo(databases)
		favService := service.NewFavouriteService(favRepo)
		favHandler := handler.NewFavouriteHandler(favService)

		favourite := v1.Group("/favourite")
		{
			favourite.POST("", favHandler.AddFavourite)
			favourite.DELETE("", favHandler.RemoveFavourite)
			favourite.GET("", favHandler.GetFavourites)
			favourite.GET("/all", favHandler.GetAllFavourites)
		}

		// Collection routes
		collectionRepo := repository.NewCollectionRepo(databases)
		collectionService := service.NewCollectionService(collectionRepo)
		collectionHandler := handler.NewCollectionHandler(collectionService)

		collection := v1.Group("/collection")
		{
			collection.POST("", collectionHandler.AddCollection)
			collection.DELETE("", collectionHandler.RemoveCollection)
			collection.GET("", collectionHandler.GetCollections)
			collection.GET("/all", collectionHandler.GetAllCollections)
		}

		// Anime routes
		animeRepo := repository.NewAnimeRepo(databases)
		animeService := service.NewAnimeService(animeRepo)
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
