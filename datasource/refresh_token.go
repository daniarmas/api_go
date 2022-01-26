package datasource

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenDatasource interface {
	GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error)
}

type refreshTokenDatasource struct{}

func (v *refreshTokenDatasource) CreateRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken) (*models.RefreshToken, error) {
	result := tx.Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshToken, nil
}

func (r *refreshTokenDatasource) DeleteRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*[]models.RefreshToken, error) {
	var refreshTokenResultSlice *[]models.RefreshToken
	result := tx.Clauses(clause.Returning{}).Where(refreshToken).Delete(&refreshTokenResultSlice)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshTokenResultSlice, nil
}

func (r *refreshTokenDatasource) GetRefreshToken(tx *gorm.DB, refreshToken *models.RefreshToken, fields *[]string) (*models.RefreshToken, error) {
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
