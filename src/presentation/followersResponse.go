package presentation

import (
	"mybook-api/src/models"
)

type Followers struct {
	UserID    string   `json:"userid"`
	Followers []string `json:"followers"`
}

type Following struct {
	UserID    string   `json:"userid"`
	Followers []string `json:"following"`
}

func NewFollowersResponse(followers models.Seguidores) Followers {
	return Followers{UserID: followers.UserID, Followers: followers.Followers}
}

func NewFollowingResponse(following models.Seguindo) Following {
	return Following{UserID: following.UserID, Followers: following.Following}
}
