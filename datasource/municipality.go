package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type MunicipalityDatasource interface {
	MunicipalityByCoordinate(tx *gorm.DB, coordinate ewkb.Point) (*models.Municipality, error)
}

type municipalityDatasource struct{}

func (i *municipalityDatasource) MunicipalityByCoordinate(tx *gorm.DB, coordinate ewkb.Point) (*models.Municipality, error) {
	var municipalityResult *models.Municipality
	p := fmt.Sprintf("'POINT(%v %v)'", coordinate.Point.Coords()[1], coordinate.Point.Coords()[0])
	result := tx.Where("ST_Contains(muncipality.polygon, ST_GeomFromText(%s, 4326)) as is_contained", p).Take(&municipalityResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return municipalityResult, nil
}
