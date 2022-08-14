package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityRepository interface {
	GetUnionBusinessAndMunicipality(tx *gorm.DB, where *entity.UnionBusinessAndMunicipality) (*entity.UnionBusinessAndMunicipality, error)
	BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*entity.UnionBusinessAndMunicipality) ([]*entity.UnionBusinessAndMunicipality, error)
	ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]entity.UnionBusinessAndMunicipalityWithMunicipality, error)
}

type unionBusinessAndMunicipality struct{}

func (v *unionBusinessAndMunicipality) GetUnionBusinessAndMunicipality(tx *gorm.DB, where *entity.UnionBusinessAndMunicipality) (*entity.UnionBusinessAndMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().GetUnionBusinessAndMunicipality(tx, &entity.UnionBusinessAndMunicipality{MunicipalityId: where.MunicipalityId, BusinessId: where.BusinessId})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *unionBusinessAndMunicipality) BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*entity.UnionBusinessAndMunicipality) ([]*entity.UnionBusinessAndMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().BatchCreateUnionBusinessAndMunicipality(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *unionBusinessAndMunicipality) ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]entity.UnionBusinessAndMunicipalityWithMunicipality, error) {
	res, err := Datasource.NewUnionBusinessAndMunicipalityDatasource().ListUnionBusinessAndMunicipalityWithMunicipality(tx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
