package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"-" bson:"password"`
	Email        string             `json:"email" bson:"email"`
	ProfileImage string             `json:"profileImage" bson:"profileImage"`
	DisplayName  string             `json:"displayName" bson:"displayName"`
}

type UserViewByOthers struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Username       string             `json:"username" bson:"username"`
	ProfileImage   string             `json:"profileImage" bson:"profileImage"`
	DisplayName    string             `json:"displayName" bson:"displayName"`
	IsFollowed     bool               `json:"isFollowed" bson:"isFollowed"`
	TotalFollowers int                `json:"totalFollowers" bson:"totalFollowers"`
	TotalFollowing int                `json:"totalFollowing" bson:"totalFollowing"`
}
