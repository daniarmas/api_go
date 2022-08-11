package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessCategoryRepository interface {
	GetBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory) (*entity.BusinessCategory, error)
	CreateBusinessCategory(tx *gorm.DB, data *entity.BusinessCategory) (*entity.BusinessCategory, error)
	DeleteBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, ids *[]uuid.UUID) (*[]entity.BusinessCategory, error)
	ListBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory) (*[]entity.BusinessCategory, error)
}

type businessCategoryRepository struct{}

func (i *businessCategoryRepository) ListBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory) (*[]entity.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().ListBusinessCategory(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessCategoryRepository) CreateBusinessCategory(tx *gorm.DB, data *entity.BusinessCategory) (*entity.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().CreateBusinessCategory(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessCategoryRepository) GetBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory) (*entity.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().GetBusinessCategory(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessCategoryRepository) DeleteBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, ids *[]uuid.UUID) (*[]entity.BusinessCategory, error) {
	res, err := Datasource.NewBusinessCategoryDatasource().DeleteBusinessCategory(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
