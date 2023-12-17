package dto

import "time"

type GetCommentByPostIDResponse struct {
	Content         string     `json:"content"`
	CreatedDatetime *time.Time `json:"createdDatetime"`
	UserID          string     `json:"userID"`
	Username        string     `json:"username"`
	DisplayName     string     `json:"displayName"`
	ProfileImage    string     `json:"profileImage"`
	CommentID       string     `json:"commentID"`
}
