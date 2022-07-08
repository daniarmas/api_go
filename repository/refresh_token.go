package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	GetRefreshToken(tx *gorm.DB, where *entity.RefreshToken, fields *[]string) (*entity.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
	DeleteRefreshTokenDeviceIdNotEqual(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
}

type refreshTokenRepository struct{}

func (v *refreshTokenRepository) CreateRefreshToken(tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error) {
	res, err := Datasource.NewRefreshTokenDatasource().CreateRefreshToken(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	res, err := Datasource.NewRefreshTokenDatasource().DeleteRefreshToken(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenDeviceIdNotEqual(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	res, err := Datasource.NewRefreshTokenDatasource().DeleteRefreshTokenDeviceIdNotEqual(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *refreshTokenRepository) GetRefreshToken(tx *gorm.DB, where *entity.RefreshToken, fields *[]string) (*entity.RefreshToken, error) {
	res, err := Datasource.NewRefreshTokenDatasource().GetRefreshToken(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
