package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenQuery interface {
	GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error)
}

type refreshTokenQuery struct{}

func (v *refreshTokenQuery) CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error) {
	result := tx.Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshToken, nil
}

func (r *refreshTokenQuery) DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error) {
	var refreshTokenResultSlice *[]models.RefreshToken
	result := tx.Clauses(clause.Returning{}).Where(refreshToken).Delete(&refreshTokenResultSlice)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshTokenResultSlice, nil
}

func (r *refreshTokenQuery) GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error) {
	var refreshTokenResult *models.RefreshToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Limit(1).Where(refreshToken).Select(*fields).Find(&refreshTokenResult)
	} else {
		result = tx.Limit(1).Where(refreshToken).Find(&refreshTokenResult)
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return refreshTokenResult, nil
		} else {
			return nil, result.Error
		}
	}
	return refreshTokenResult, nil
}
