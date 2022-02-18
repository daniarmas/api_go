package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UnionBusinessAndMunicipalityDatasource interface {
	UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error
}

type unionBusinessAndMunicipalityDatasource struct{}

func (v *unionBusinessAndMunicipalityDatasource) UnionBusinessAndMunicipalityExists(tx *gorm.DB, where *models.UnionBusinessAndMunicipality) error {
	var res models.UnionBusinessAndMunicipality
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return errors.New("record not found")
		} else {
			return result.Error
		}
	}
	return nil
}
