package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityRepository interface {
	GetUnionBusinessAndMunicipality(tx *gorm.DB, where *models.UnionBusinessAndMunicipality, fields *[]string) (*models.UnionBusinessAndMunicipality, error)
	BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*models.UnionBusinessAndMunicipality) ([]*models.UnionBusinessAndMunicipality, error)
	ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]models.UnionBusinessAndMunicipalityWithMunicipality, error)
}

type unionBusinessAndMunicipality struct{}

func (v *unionBusinessAndMunicipality) GetUnionBusinessAndMunicipality(tx *gorm.DB, where *models.UnionBusinessAndMunicipality, fields *[]string) (*models.UnionBusinessAndMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().GetUnionBusinessAndMunicipality(tx, &models.UnionBusinessAndMunicipality{MunicipalityId: where.MunicipalityId, BusinessId: where.BusinessId}, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
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
