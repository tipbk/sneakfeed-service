package dto

import "github.com/tipbk/sneakfeed-service/model"

type GetUserByUsernameResponse struct {
	IsYourUser bool `json:"isYourUser"`
	*model.UserViewByOthers
}
