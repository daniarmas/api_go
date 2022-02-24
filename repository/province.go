package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetProvince(tx *gorm.DB, where *models.Province) (*models.Province, error)
}

type provinceRepository struct{}

func (v *provinceRepository) GetProvince(tx *gorm.DB, where *models.Province) (*models.Province, error) {
	result, err := Datasource.NewProvinceDatasource().GetProvince(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
