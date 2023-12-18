package dto

import "github.com/tipbk/sneakfeed-service/model"

type GetUserByUsernameResponse struct {
	IsFollowed bool        `json:"isFollowed"`
	IsYourUser bool        `json:"isYourUser"`
	User       *model.User `json:"user"`
}
