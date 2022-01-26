package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type RefreshTokenQuery interface {
	GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error)
}

type refreshTokenQuery struct{}

func (v *refreshTokenQuery) CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error) {
	result, err := Datasource.NewRefreshTokenDatasource().CreateRefreshToken(tx, refreshToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *refreshTokenQuery) DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error) {
	result, err := Datasource.NewRefreshTokenDatasource().DeleteRefreshToken(tx, refreshToken, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *refreshTokenQuery) GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error) {
	result, err := Datasource.NewRefreshTokenDatasource().GetRefreshToken(tx, refreshToken, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
