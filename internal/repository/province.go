package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetProvince(tx *gorm.DB, where *entity.Province) (*entity.Province, error)
}

type provinceRepository struct{}

func (v *provinceRepository) GetProvince(tx *gorm.DB, where *entity.Province) (*entity.Province, error) {
	res, err := Datasource.NewProvinceDatasource().GetProvince(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
