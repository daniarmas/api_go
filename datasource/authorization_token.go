package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthorizationTokenDatasource interface {
	GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error)
	DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error)
}

type authorizationTokenDatasource struct{}

func (v *authorizationTokenDatasource) CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	result := tx.Create(&authorizationToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationToken, nil
}

func (r *authorizationTokenDatasource) DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error) {
	var authorizationTokenResult *[]models.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(where, *ids).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (r *authorizationTokenDatasource) DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error) {
	var authorizationTokenResult *[]models.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(authorizationToken).Delete(&authorizationTokenResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return authorizationTokenResult, nil
}

func (v *authorizationTokenDatasource) GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error) {
	var authorizationTokenResult *models.AuthorizationToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Where(authorizationToken).Select(*fields).Take(&authorizationTokenResult)
	} else {
		result = tx.Where(authorizationToken).Take(&authorizationTokenResult)
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("authorizationtoken not found")
		} else {
			return nil, result.Error
		}
	}
	return authorizationTokenResult, nil
}
