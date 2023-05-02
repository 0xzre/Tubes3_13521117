package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	ID       primitive.ObjectID `bson:"_id"`
	Question *string            `json:"question"`
	Answer   *string            `json:"answer"`
}
