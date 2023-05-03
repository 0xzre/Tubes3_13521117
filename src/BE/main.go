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

	// these are the endpoints
	//C
	router.POST("/question/create", controller.AddQuestion)
	//R
	router.GET("/answer/KMP/:question", controller.GetResponseKMP)
	// router.GET("/answer/BM/:question", controller.GetResponseBM)
	//U
	router.PUT("/answer/update/:question", controller.UpdateAnswer)
	router.PUT("/question/update/:answer", controller.UpdateQuestion)
	//D
	router.DELETE("/question/delete/:question", controller.DeleteQuestion)

	//this runs the server and allows it to listen to requests.
	router.Run("localhost:" + port)
}
