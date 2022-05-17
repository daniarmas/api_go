package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityDatasource interface {
	GetUnionBusinessAndMunicipality(tx *gorm.DB, where *models.UnionBusinessAndMunicipality, fields *[]string) (*models.UnionBusinessAndMunicipality, error)
	BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*models.UnionBusinessAndMunicipality) ([]*models.UnionBusinessAndMunicipality, error)
	ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]models.UnionBusinessAndMunicipalityWithMunicipality, error)
}

type unionBusinessAndMunicipalityDatasource struct{}

func (v *unionBusinessAndMunicipalityDatasource) GetUnionBusinessAndMunicipality(tx *gorm.DB, where *models.UnionBusinessAndMunicipality, fields *[]string) (*models.UnionBusinessAndMunicipality, error) {
	var res *models.UnionBusinessAndMunicipality
	result := tx.Where(where).Select(*fields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("refreshtoken not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *unionBusinessAndMunicipalityDatasource) BatchCreateUnionBusinessAndMunicipality(tx *gorm.DB, data []*models.UnionBusinessAndMunicipality) ([]*models.UnionBusinessAndMunicipality, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *unionBusinessAndMunicipalityDatasource) ListUnionBusinessAndMunicipalityWithMunicipality(tx *gorm.DB, ids []string) (*[]models.UnionBusinessAndMunicipalityWithMunicipality, error) {
	var response []models.UnionBusinessAndMunicipalityWithMunicipality
	result := tx.Model(&models.UnionBusinessAndMunicipality{}).Select(`union_business_and_municipality.id, union_business_and_municipality.business_id, union_business_and_municipality.municipality_id, municipality.name as municipality_name, union_business_and_municipality.create_time, union_business_and_municipality.update_time `).Joins("left join municipality on municipality.id = union_business_and_municipality.municipality_id").Where("union_business_and_municipality.id IN ? ", ids).Find(&response)
	if result.Error != nil {
		return nil, result.Error
	}
	return &response, nil
}
