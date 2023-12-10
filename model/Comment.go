package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	UserID          string             `json:"userID" bson:"userID"`
	PostID          string             `json:"postID" bson:"postID"`
	Content         string             `json:"content" bson:"content"`
	CreatedDatetime *time.Time         `json:"createdDatetime" bson:"createdDatetime"`
}
