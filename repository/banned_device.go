package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedDeviceQuery interface {
	GetBannedDevice(tx *gorm.DB, where *models.BannedDevice) (*models.BannedDevice, error)
}

type bannedDeviceQuery struct{}

func (i *bannedDeviceQuery) GetBannedDevice(tx *gorm.DB, where *models.BannedDevice) (*models.BannedDevice, error) {
	res, err := Datasource.NewBannedDeviceDatasource().GetBannedDevice(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
