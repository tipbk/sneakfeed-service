package dto

import "time"

type GetPostsResponse struct {
	PostID          string     `json:"postID"`
	UserID          string     `json:"userID"`
	Username        string     `json:"username"`
	ProfileImage    string     `json:"profileImage"`
	Title           string     `json:"title"`
	CreatedDatetime *time.Time `json:"createdDatetime"`
	Content         string     `json:"content"`
	PostImageUrl    *string    `json:"postImageUrl"`
	TotalComment    int64      `json:"totalComment"`
	TotalLike       int64      `json:"totalLike"`
}
