package dto

import "time"

type GetCommentByPostIDResponse struct {
	Content         string     `json:"content"`
	CreatedDatetime *time.Time `json:"createdDatetime"`
	UserID          string     `json:"userID"`
	Username        string     `json:"username"`
	ProfileImage    string     `json:"profileImage"`
	CommentID       string     `json:"commentID"`
}
