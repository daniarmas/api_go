package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeviceDatasource interface {
	GetDevice(tx *gorm.DB, device *models.Device, fields *[]string) (*models.Device, error)
	CreateDevice(tx *gorm.DB, device *models.Device) (*models.Device, error)
	UpdateDevice(tx *gorm.DB, where *models.Device, device *models.Device) (*models.Device, error)
}

type deviceDatasource struct{}

func (i *deviceDatasource) GetDevice(tx *gorm.DB, device *models.Device, fields *[]string) (*models.Device, error) {
	var deviceResult *models.Device
	result := tx.Where(device).Select(*fields).Take(&deviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return deviceResult, nil
}

func (v *deviceDatasource) CreateDevice(tx *gorm.DB, device *models.Device) (*models.Device, error) {
	result := tx.Create(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func (v *deviceDatasource) UpdateDevice(tx *gorm.DB, where *models.Device, device *models.Device) (*models.Device, error) {
	result := tx.Where(where).Updates(&device)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return device, nil
}
