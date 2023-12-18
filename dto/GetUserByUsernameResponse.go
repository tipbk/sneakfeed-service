package dto

import "github.com/tipbk/sneakfeed-service/model"

type GetUserByUsernameResponse struct {
	IsFollowed bool        `json:"isFollowed"`
	User       *model.User `json:"user"`
}
