package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedDeviceQuery interface {
	GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error)
}

type bannedDeviceQuery struct{}

func (i *bannedDeviceQuery) GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error) {
	result, err := Datasource.NewBannedDeviceDatasource().GetBannedDevice(tx, bannedDevice, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
