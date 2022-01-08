package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type BannedDeviceQuery interface {
	GetBannedDevice(tx *gorm.DB, bannedDevice *datastruct.BannedDevice, fields *[]string) (*datastruct.BannedDevice, error)
	// ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type bannedDeviceQuery struct{}

func (i *bannedDeviceQuery) GetBannedDevice(tx *gorm.DB, bannedDevice *datastruct.BannedDevice, fields *[]string) (*datastruct.BannedDevice, error) {
	var bannedDeviceResult *datastruct.BannedDevice
	result := tx.Limit(1).Where(bannedDevice).Select(*fields).Find(&bannedDeviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return bannedDeviceResult, nil
		} else {
			return nil, result.Error
		}
	}
	return bannedDeviceResult, nil
}
