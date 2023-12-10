package dto

type CreatePostRequest struct {
	Content     string  `json:"content"`
	ImageBase64 *string `json:"imageBase64"`
}
