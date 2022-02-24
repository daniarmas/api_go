package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type ProvinceDatasource interface {
	ProvinceByCoordinate(tx *gorm.DB, coordinate ewkb.Point) (*models.Province, error)
	GetProvince(tx *gorm.DB, where *models.Province) (*models.Province, error)
}

type provinceDatasource struct{}

func (i *provinceDatasource) ProvinceByCoordinate(tx *gorm.DB, coordinate ewkb.Point) (*models.Province, error) {
	var provinceResult *models.Province
	p := fmt.Sprintf("'POINT(%v %v)'", coordinate.Point.Coords()[1], coordinate.Point.Coords()[0])
	result := tx.Where("ST_Contains(province.polygon, ST_GeomFromText(%s, 4326)) as is_contained", p).Take(&provinceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return provinceResult, nil
}

func (i *provinceDatasource) GetProvince(tx *gorm.DB, where *models.Province) (*models.Province, error) {
	var response *models.Province
	result := tx.Where(where).Take(&response)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return response, nil
}
