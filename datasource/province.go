package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type ProvinceDatasource interface {
	ProvinceByCoordinate(tx *gorm.DB, coordinate ewkb.Point, fields *[]string) (*models.Province, error)
	GetProvince(tx *gorm.DB, where *models.Province, fields *[]string) (*models.Province, error)
}

type provinceDatasource struct{}

func (i *provinceDatasource) ProvinceByCoordinate(tx *gorm.DB, coordinate ewkb.Point, fields *[]string) (*models.Province, error) {
	var provinceResult *models.Province
	p := fmt.Sprintf("'POINT(%v %v)'", coordinate.Point.Coords()[1], coordinate.Point.Coords()[0])
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where("ST_Contains(province.polygon, ST_GeomFromText(%s, 4326)) as is_contained", p).Select(*selectFields).Take(&provinceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return provinceResult, nil
}

func (i *provinceDatasource) GetProvince(tx *gorm.DB, where *models.Province, fields *[]string) (*models.Province, error) {
	var res *models.Province
	selectFields := &[]string{"*"}
	if fields == nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
