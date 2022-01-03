package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type AuthorizationTokenQuery interface {
	// GetDevice(device *datastruct.Device, fields *[]string) (*[]datastruct.Device, error)
	// ListItem() ([]datastruct.Item, error)
	CreateAuthorizationToken(authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error)
	// UpdateDevice(device *datastruct.Device) error
	// DeleteRefreshToken(refreshToken *datastruct.RefreshToken) error
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error) {
	result := DB.Table("AuthorizationToken").Create(&authorizationToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationToken, nil
}
