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
	p := fmt.Sprintf("POINT(%v %v)", coordinate.Point.Coords()[1], coordinate.Point.Coords()[0])
	result := tx.Select("id, name, province_id, ST_AsEWKB(coordinates) AS coordinates, zoom, create_time, update_time").Where("ST_Contains(polygon, ST_GeomFromText(?, 4326))", p).Take(&municipalityResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return municipalityResult, nil
}
