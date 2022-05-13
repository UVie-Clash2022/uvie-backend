package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username" validate:"required,min=2,max=100"`
	Password string             `json:"password,omitempty" bson:"password" validate:"required,min=4,max=100"`
}
