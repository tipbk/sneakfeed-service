package model

type Pagination struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
}
