package main

import (
	"os"

	"BE/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	// Endpoints used
	router.GET("/answer/KMP/:question", controller.GetResponseKMP)
	// router.GET("/answer/BM/:question", controller.GetResponseBM)

	// Runs the server and allows it to listen to requests
	router.Run("localhost:" + port)
}
