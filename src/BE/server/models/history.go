package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type History struct {
	ID       primitive.ObjectID `bson:"_id"`
	Prompt   *string            `json:"prompt"`
	Response *string            `json:"response"`
}
