package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type BannedDeviceQuery interface {
	GetBannedDevice(bannedDevice *datastruct.BannedDevice) (*[]datastruct.BannedDevice, error)
	// ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type bannedDeviceQuery struct{}

func (i *bannedDeviceQuery) GetBannedDevice(bannedDevice *datastruct.BannedDevice) (*[]datastruct.BannedDevice, error) {
	var bannedDeviceResult *[]datastruct.BannedDevice
	result := DB.Table("BannedDevice").Limit(1).Where(bannedDevice).Find(&bannedDeviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return bannedDeviceResult, nil
		} else {
			return nil, result.Error
		}
	}
	return bannedDeviceResult, nil
}
