package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorizationTokenQuery interface {
	GetAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken) (*models.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, data *models.AuthorizationToken) (*models.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error)
	DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error)
}

type authorizationTokenQuery struct{}

func (v *authorizationTokenQuery) CreateAuthorizationToken(tx *gorm.DB, data *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	res, err := Datasource.NewAuthorizationTokenDatasource().CreateAuthorizationToken(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error) {
	res, err := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationToken(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *authorizationTokenQuery) DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error) {
	res, err := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationTokenByRefreshTokenIds(tx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *authorizationTokenQuery) GetAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	res, err := Datasource.NewAuthorizationTokenDatasource().GetAuthorizationToken(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
