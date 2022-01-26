package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type AuthorizationTokenQuery interface {
	GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error)
	DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error)
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	result, err := Datasource.NewAuthorizationTokenDatasource().CreateAuthorizationToken(tx, authorizationToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationTokenIn(tx *gorm.DB, where string, ids *[]string) (*[]models.AuthorizationToken, error) {
	result, err := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationTokenIn(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken) (*[]models.AuthorizationToken, error) {
	result, err := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationToken(tx, authorizationToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *authorizationTokenQuery) GetAuthorizationToken(tx *gorm.DB, authorizationToken *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error) {
	result, err := Datasource.NewAuthorizationTokenDatasource().GetAuthorizationToken(tx, authorizationToken, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
