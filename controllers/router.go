package controllers

import (
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		gamesGroup := apiGroup.Group("/game")
		{
			gamesGroup.POST("", StartGame)

			gameGroup := gamesGroup.Group("/:gameId")
			{
				gameGroup.GET("", GetGame)
				gameGroup.POST("/move", MakeMove)
			}
		}
	}

	return router
}
