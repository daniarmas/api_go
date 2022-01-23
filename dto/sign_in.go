package dto

import "github.com/daniarmas/api_go/models"

type SignIn struct {
	RefreshToken       string
	AuthorizationToken string
	User               models.User
}
