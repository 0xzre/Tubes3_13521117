package main

import (
	"BE/controllers/questionController"
	"BE/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	r.GET("/api/questions", questionController.IndexQuestion)
	r.GET("/api/question/:pertanyaan", questionController.ShowQuestion)
	r.POST("/api/question", questionController.CreateQuestion)
	r.PUT("/api/question/:pertanyaan", questionController.UpdateQuestion)
	r.DELETE("/api/question", questionController.DeleteQuestion)

	r.Run()
}
