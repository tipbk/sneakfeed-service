package dto

type UpdateUserProfileRequest struct {
	ImageBase64 string `json:"imageBase64"`
	DisplayName string `json:"displayName"`
}
