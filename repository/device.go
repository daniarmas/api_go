package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDevice(tx *gorm.DB, where *models.Device, fields *[]string) (*models.Device, error)
	CreateDevice(tx *gorm.DB, data *models.Device) (*models.Device, error)
	UpdateDevice(tx *gorm.DB, where *models.Device, data *models.Device) (*models.Device, error)
}

type deviceRepository struct{}

func (i *deviceRepository) GetDevice(tx *gorm.DB, where *models.Device, fields *[]string) (*models.Device, error) {
	res, err := Datasource.NewDeviceDatasource().GetDevice(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (v *deviceRepository) CreateDevice(tx *gorm.DB, data *models.Device) (*models.Device, error) {
	res, err := Datasource.NewDeviceDatasource().CreateDevice(tx, data)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (v *deviceRepository) UpdateDevice(tx *gorm.DB, where *models.Device, data *models.Device) (*models.Device, error) {
	res, err := Datasource.NewDeviceDatasource().UpdateDevice(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, err
}
