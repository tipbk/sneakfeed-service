package dto

import "time"

type GetPostByIDResponse struct {
	PostID          string     `json:"postID"`
	UserID          string     `json:"userID"`
	Username        string     `json:"username"`
	ProfileImage    string     `json:"profileImage"`
	Content         string     `json:"content"`
	CreatedDatetime *time.Time `json:"createdDatetime"`
	IsLike          bool       `json:"isLike"`
	PostImageUrl    *string    `json:"postImageUrl"`
	TotalLikes      int64      `json:"totalLikes"`
	TotalComments   int64      `json:"totalComments"`
	IsComment       bool       `json:"isComment"`
}
