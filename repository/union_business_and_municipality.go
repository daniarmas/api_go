package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityRepository interface {
	UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error
	BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*models.UnionBusinessAndMunicipality) ([]*models.UnionBusinessAndMunicipality, error)
	ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]models.UnionBusinessAndMunicipalityWithMunicipality, error)
}

type unionBusinessAndMunicipality struct{}

func (v *unionBusinessAndMunicipality) UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error {
	err := Datasource.NewUnionBusinessAndMunicipalityDatasource().UnionBusinessAndMunicipalityExists(tx, &models.UnionBusinessAndMunicipality{MunicipalityFk: where.MunicipalityFk, BusinessFk: where.BusinessFk})
	if err != nil {
		return err
	}
	return nil
}

func (v *unionBusinessAndMunicipality) BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*models.UnionBusinessAndMunicipality) ([]*models.UnionBusinessAndMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().BatchCreateUnionBusinessAndMunicipality(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *unionBusinessAndMunicipality) ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]models.UnionBusinessAndMunicipalityWithMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().ListUnionBusinessAndMunicipalityWithMunicipality(tx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
