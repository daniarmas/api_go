package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type MunicipalityRepository interface {
	GetMunicipalityByCoordinate(tx *gorm.DB, coordinates ewkb.Point) (*models.Municipality, error)
}

type municipalityRepository struct{}

func (v *municipalityRepository) GetMunicipalityByCoordinate(tx *gorm.DB, coordinates ewkb.Point) (*models.Municipality, error) {
	result, err := Datasource.NewMunicipalityDatasource().MunicipalityByCoordinate(tx, coordinates)
	if err != nil {
		return nil, err
	}
	return result, nil
}
