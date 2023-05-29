package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"jakebaum.uk/tictactoe/controllers"
)

func main() {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		gamesGroup := apiGroup.Group("/game")
		{
			gamesGroup.POST("", controllers.StartGame)

			gameGroup := gamesGroup.Group("/:gameId")
			{
				gameGroup.GET("", controllers.GetGame)
				gameGroup.POST("/move", controllers.MakeMove)
			}
		}
	}

	err := router.Run("localhost:8080")
	if err != nil {
		log.Error("An error occurred when trying to start server: ", err)
		return
	}
}
