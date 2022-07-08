package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type DeviceDatasource interface {
	GetDevice(tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error)
	CreateDevice(tx *gorm.DB, data *entity.Device) (*entity.Device, error)
	UpdateDevice(tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error)
}

type deviceDatasource struct{}

func (i *deviceDatasource) GetDevice(tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error) {
	var res *entity.Device
	selectFields := &[]string{"*"}
	if fields != nil {
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

func (v *deviceDatasource) CreateDevice(tx *gorm.DB, data *entity.Device) (*entity.Device, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *deviceDatasource) UpdateDevice(tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error) {
	result := tx.Where(where).Updates(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}
