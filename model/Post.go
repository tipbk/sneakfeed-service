package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	UserID          string             `json:"userID" bson:"userID"`
	Content         string             `json:"content" bson:"content"`
	CreatedDatetime *time.Time         `json:"createdDatetime" bson:"createdDatetime"`
	ImageUrl        *string            `json:"imageUrl" bson:"imageUrl"`
}

type PostDetail struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	UserID          string             `json:"userID" bson:"userID"`
	Content         string             `json:"content" bson:"content"`
	CreatedDatetime *time.Time         `json:"createdDatetime" bson:"createdDatetime"`
	ImageUrl        *string            `json:"imageUrl" bson:"imageUrl"`
	Username        string             `json:"username"`
	ProfileImage    string             `json:"profileImage"`
	TotalLikes      int                `json:"totalLikes"`
	TotalComments   int                `json:"totalComments"`
}
