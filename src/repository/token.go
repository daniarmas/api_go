package repository

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type TokenQuery interface {
	CreateJwtRefreshToken(refreshTokenFk *string) (*string, error)
	CreateJwtAuthorizationToken(authorizationTokenFk *string) (*string, error)
}

type tokenQuery struct{}

func (v *tokenQuery) CreateJwtRefreshToken(refreshTokenFk *string) (*string, error) {
	hmacSampleSecret := viper.Get("JWT_SECRET").(string)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refreshTokenFk": *refreshTokenFk,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (r *tokenQuery) CreateJwtAuthorizationToken(authorizationTokenFk *string) (*string, error) {
	hmacSampleSecret := viper.Get("JWT_SECRET").(string)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorizationTokenFk": *authorizationTokenFk,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
