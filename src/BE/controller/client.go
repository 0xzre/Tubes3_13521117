package controller

import (
	"BE/server/routes"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

// Client Database instance
var validate = validator.New()
var Client *mongo.Client = routes.DBinstance()
var questionCollection *mongo.Collection = routes.OpenCollection(Client, "questions")
