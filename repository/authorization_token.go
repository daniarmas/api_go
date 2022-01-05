package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type AuthorizationTokenQuery interface {
	GetAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken, fields *[]string) (*datastruct.AuthorizationToken, error)
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

func (v *authorizationTokenQuery) GetAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken, fields *[]string) (*datastruct.AuthorizationToken, error) {
	var authorizationTokenResult *datastruct.AuthorizationToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Table("AuthorizationToken").Limit(1).Where(authorizationToken).Select(*fields).Find(&authorizationTokenResult)
	} else {
		result = tx.Table("AuthorizationToken").Limit(1).Where(authorizationToken).Find(&authorizationTokenResult)
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return authorizationTokenResult, nil
		} else {
			return nil, result.Error
		}
	}
	return authorizationTokenResult, nil
}
