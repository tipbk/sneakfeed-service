package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type LikePost struct {
	ID     primitive.ObjectID `json:"-" bson:"_id"`
	UserID string             `json:"userID" bson:"userID"`
	PostID string             `json:"postID" bson:"postID"`
}
