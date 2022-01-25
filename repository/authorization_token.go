package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthorizationTokenQuery interface {
	GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error)
	DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error)
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	result := tx.Create(&authorizationToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationToken, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error) {
	var authorizationTokenResult *[]models.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(where, *ids).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error) {
	var authorizationTokenResult *[]models.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(authorizationToken).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (v *authorizationTokenQuery) GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error) {
	var authorizationTokenResult *models.AuthorizationToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Where(authorizationToken).Select(*fields).Take(&authorizationTokenResult)
	} else {
		result = tx.Where(authorizationToken).Take(&authorizationTokenResult)
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
