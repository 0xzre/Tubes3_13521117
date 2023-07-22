package main

import (
	"BE/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve static files from the "./web" directory for the root path "/"
	r.Use(static.Serve("/", static.LocalFile("./web", true)))

	// Create a new sub-router for the "/response" group
	router := r.Group("/response")

	// Use Gin Logger middleware for logging
	router.Use(gin.Logger())

	// Set up CORS middleware to allow requests from the frontend "https://cuakgpt.vercel.app"
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://cuakgpt.vercel.app"},
	}))

	// Register three endpoints with corresponding handler functions from the controller package
	router.GET("/KMP/:question", controller.GetResponseKMP)
	router.GET("/BM/:question", controller.GetResponseBM)
	router.GET("/history", controller.GetAllHistory)

	// Runs the server and allows it to listen to incoming requests on localhost:5000
	r.Run(":8080")
}
