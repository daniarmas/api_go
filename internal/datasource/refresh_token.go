package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenDatasource interface {
	GetRefreshToken(tx *gorm.DB, where *entity.RefreshToken) (*entity.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error)
	DeleteRefreshToken(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
	DeleteRefreshTokenDeviceIdNotEqual(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
}

type refreshTokenDatasource struct{}

func (v *refreshTokenDatasource) CreateRefreshToken(tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *refreshTokenDatasource) DeleteRefreshToken(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	var res *[]entity.RefreshToken
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

func (r *refreshTokenDatasource) DeleteRefreshTokenDeviceIdNotEqual(tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	var res *[]entity.RefreshToken
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(`device_id != ?`, where.DeviceId).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *refreshTokenDatasource) GetRefreshToken(tx *gorm.DB, where *entity.RefreshToken) (*entity.RefreshToken, error) {
	var res *entity.RefreshToken
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("refreshtoken not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
