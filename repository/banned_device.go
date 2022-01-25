package repository

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedDeviceQuery interface {
	GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error)
	// ListItem() ([]models.Item, error)
	// CreateItem(answer models.Item) (*int64, error)
	// UpdateItem(answer models.Item) (*models.Item, error)
	// DeleteItem(id int64) error
}

type bannedDeviceQuery struct{}

func (i *bannedDeviceQuery) GetBannedDevice(tx *gorm.DB, bannedDevice *models.BannedDevice, fields *[]string) (*models.BannedDevice, error) {
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
