package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeviceQuery interface {
	GetDevice(tx *gorm.DB, device *models.Device, fields *[]string) (*models.Device, error)
	CreateDevice(tx *gorm.DB, device *models.Device) (*models.Device, error)
	UpdateDevice(tx *gorm.DB, where *models.Device, device *models.Device) (*models.Device, error)
}

type deviceQuery struct{}

func (i *deviceQuery) GetDevice(tx *gorm.DB, where *models.Device, fields *[]string) (*models.Device, error) {
	result, err := Datasource.NewDeviceDatasource().GetDevice(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (v *deviceQuery) CreateDevice(tx *gorm.DB, data *models.Device) (*models.Device, error) {
	result, err := Datasource.NewDeviceDatasource().CreateDevice(tx, data)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (v *deviceQuery) UpdateDevice(tx *gorm.DB, where *models.Device, data *models.Device) (*models.Device, error) {
	result, err := Datasource.NewDeviceDatasource().UpdateDevice(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, err
}
