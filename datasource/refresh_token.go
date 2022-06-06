package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenDatasource interface {
	GetRefreshToken(tx *gorm.DB, where *models.RefreshToken, fields *[]string) (*models.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, data *models.RefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, where *models.RefreshToken, ids *[]uuid.UUID) (*[]models.RefreshToken, error)
}

type refreshTokenDatasource struct{}

func (v *refreshTokenDatasource) CreateRefreshToken(tx *gorm.DB, data *models.RefreshToken) (*models.RefreshToken, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *refreshTokenDatasource) DeleteRefreshToken(tx *gorm.DB, where *models.RefreshToken, ids *[]uuid.UUID) (*[]models.RefreshToken, error) {
	var res *[]models.RefreshToken
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *refreshTokenDatasource) GetRefreshToken(tx *gorm.DB, where *models.RefreshToken, fields *[]string) (*models.RefreshToken, error) {
	var res *models.RefreshToken
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("refreshtoken not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
