package dto

import "github.com/daniarmas/api_go/datastruct"

type SignIn struct {
	RefreshToken       string
	AuthorizationToken string
	User               datastruct.User
}
