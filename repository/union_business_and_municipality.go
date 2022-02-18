package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityRepository interface {
	UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error
}

type unionBusinessAndMunicipality struct{}

func (v *unionBusinessAndMunicipality) UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error {
	err := Datasource.NewUnionBusinessAndMunicipalityDatasource().UnionBusinessAndMunicipalityExists(tx, &models.UnionBusinessAndMunicipality{MunicipalityFk: where.MunicipalityFk, BusinessFk: where.BusinessFk})
	if err != nil {
		return err
	}
	return nil
}
