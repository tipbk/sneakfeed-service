package dto

type CreatePostRequest struct {
	Content       string  `json:"content"`
	ImageBase64   *string `json:"imageBase64"`
	OgTitle       *string `json:"ogTitle"`
	OgDescription *string `json:"ogDescription"`
	OgLink        *string `json:"ogLink"`
	OgImage       *string `json:"ogImage"`
	OgDomain      *string `json:"ogDomain"`
}
