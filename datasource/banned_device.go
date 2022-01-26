package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedDeviceDatasource interface {
	GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error)
}

type bannedDeviceDatasource struct{}

func (i *bannedDeviceDatasource) GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error) {
	var bannedDeviceResult *models.BannedDevice
	result := tx.Where(bannedDevice).Select(*fields).Take(&bannedDeviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return bannedDeviceResult, nil
}
