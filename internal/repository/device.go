package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDevice(tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error)
	CreateDevice(tx *gorm.DB, data *entity.Device) (*entity.Device, error)
	UpdateDevice(tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error)
}

type deviceRepository struct{}

func (i *deviceRepository) GetDevice(tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error) {
	res, err := Datasource.NewDeviceDatasource().GetDevice(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (v *deviceRepository) CreateDevice(tx *gorm.DB, data *entity.Device) (*entity.Device, error) {
	res, err := Datasource.NewDeviceDatasource().CreateDevice(tx, data)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (v *deviceRepository) UpdateDevice(tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error) {
	res, err := Datasource.NewDeviceDatasource().UpdateDevice(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, err
}
