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
	OgTitle         *string            `json:"ogTitle" bson:"ogTitle"`
	OgDescription   *string            `json:"ogDescription" bson:"ogDescription"`
	OgLink          *string            `json:"ogLink" bson:"ogLink"`
	OgImage         *string            `json:"ogImage" bson:"ogImage"`
	OgDomain        *string            `json:"ogDoamin" bson:"ogDomain"`
}

type PostDetail struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	UserID          string             `json:"userID" bson:"userID"`
	Content         string             `json:"content" bson:"content"`
	CreatedDatetime *time.Time         `json:"createdDatetime" bson:"createdDatetime"`
	ImageUrl        *string            `json:"imageUrl" bson:"imageUrl"`
	Username        string             `json:"username"`
	DisplayName     string             `json:"displayName"`
	ProfileImage    string             `json:"profileImage"`
	TotalLikes      int                `json:"totalLikes"`
	TotalComments   int                `json:"totalComments"`
	IsLike          bool               `json:"isLike"`
	IsComment       bool               `json:"isComment"`
}
type PostDetailPagination struct {
	Pagination Pagination   `json:"pagination" bson:"pagination"`
	Posts      []PostDetail `json:"posts" bson:"posts"`
}
