package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follow struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	UserID          string             `json:"userID" bson:"userID"`
	FollowUserID    string             `json:"followUserID" bson:"followUserID"`
	CreatedDatetime *time.Time         `json:"createdDatetime" bson:"createdDatetime"`
}
