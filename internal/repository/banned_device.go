package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BannedDeviceRepository interface {
	GetBannedDevice(tx *gorm.DB, where *entity.BannedDevice, fields *[]string) (*entity.BannedDevice, error)
}

type bannedDeviceRepository struct{}

func (i *bannedDeviceRepository) GetBannedDevice(tx *gorm.DB, where *entity.BannedDevice, fields *[]string) (*entity.BannedDevice, error) {
	res, err := Datasource.NewBannedDeviceDatasource().GetBannedDevice(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
