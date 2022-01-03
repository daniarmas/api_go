package dto

import "github.com/daniarmas/api_go/src/datastruct"

type SignIn struct {
	RefreshToken       string
	AuthorizationToken string
	User               datastruct.User
}
