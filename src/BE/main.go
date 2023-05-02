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
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(cors.Default())

	// these are the endpoints
	//C
	router.POST("/question/create", controller.AddQuestion)
	//R
	router.GET("/answer/:answer", controller.GetAnswerByQuestion)
	router.GET("/questions", controller.GetQuestions)
	router.GET("/question/:id/", controller.GetQuestionById)
	//U
	router.PUT("/answer/update/:id", controller.UpdateAnswer)
	router.PUT("/question/update/:id", controller.UpdateQuestion)
	//D
	router.DELETE("/question/delete/:id", controller.DeleteQuestion)

	//this runs the server and allows it to listen to requests.
	router.Run("localhost:" + port)
}
