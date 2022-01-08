package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthorizationTokenQuery interface {
	GetAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken, fields *[]string) (*datastruct.AuthorizationToken, error)
	// ListItem() ([]datastruct.Item, error)
	CreateAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error)
	// UpdateDevice(device *datastruct.Device) error
	DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*[]datastruct.AuthorizationToken, error)
	DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]datastruct.AuthorizationToken, error)
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*datastruct.AuthorizationToken, error) {
	result := tx.Create(&authorizationToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationToken, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]datastruct.AuthorizationToken, error) {
	var authorizationTokenResult *[]datastruct.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(where, *ids).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken) (*[]datastruct.AuthorizationToken, error) {
	var authorizationTokenResult *[]datastruct.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(authorizationToken).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (v *authorizationTokenQuery) GetAuthorizationToken(tx *gorm.DB, authorizationToken *datastruct.AuthorizationToken, fields *[]string) (*datastruct.AuthorizationToken, error) {
	var authorizationTokenResult *datastruct.AuthorizationToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Limit(1).Where(authorizationToken).Select(*fields).Find(&authorizationTokenResult)
	} else {
		result = tx.Limit(1).Where(authorizationToken).Find(&authorizationTokenResult)
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
