package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"-" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"-" bson:"password"`
	ProfileImage string             `json:"profileImage" bson:"profileImage"`
}
