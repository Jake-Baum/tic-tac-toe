package main

import (
	log "github.com/sirupsen/logrus"
	"jakebaum.uk/tictactoe/controllers"
)

func main() {
	router := controllers.SetUpRouter()

	err := router.Run("localhost:8080")
	if err != nil {
		log.Error("An error occurred when trying to start server: ", err)
		return
	}
}
