package main

import (
	"BE/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Dont worry about this line just yet, it will make sense in the Dockerise bit!
	r.Use(static.Serve("/", static.LocalFile("./web", true)))
	router := r.Group("/response")

	router.Use(gin.Logger())
	router.Use(cors.Default())

	// Endpoints used
	router.GET("/KMP/:question", controller.GetResponseKMP)
	// router.GET("/BM/:question", controller.GetResponseBM)

	// Runs the server and allows it to listen to requests
	// Runs in localhost 5000
	r.Run()
}
