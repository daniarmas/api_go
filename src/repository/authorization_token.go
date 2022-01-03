package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"gorm.io/gorm"
)

type AuthorizationTokenQuery interface {
	// GetDevice(device *datastruct.Device, fields *[]string) (*[]datastruct.Device, error)
	// ListItem() ([]datastruct.Item, error)
	CreateAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error)
	// UpdateDevice(device *datastruct.Device) error
	// DeleteRefreshToken(refreshToken *datastruct.RefreshToken) error
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error) {
	result := tx.Table("AuthorizationToken").Create(&authorizationToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationToken, nil
}
