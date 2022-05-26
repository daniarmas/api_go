package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetProvince(tx *gorm.DB, where *models.Province, fields *[]string) (*models.Province, error)
}

type provinceRepository struct{}

func (v *provinceRepository) GetProvince(tx *gorm.DB, where *models.Province, fields *[]string) (*models.Province, error) {
	res, err := Datasource.NewProvinceDatasource().GetProvince(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
