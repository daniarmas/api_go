package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type DeviceQuery interface {
	GetDevice(device *datastruct.Device, fields *[]string) (*datastruct.Device, error)
	// ListItem() ([]datastruct.Item, error)
	CreateDevice(device *datastruct.Device) (*datastruct.Device, error)
	UpdateDevice(where *datastruct.Device, device *datastruct.Device) (*datastruct.Device, error)
	// DeleteItem(id int64) error
}

type deviceQuery struct{}

func (i *deviceQuery) GetDevice(device *datastruct.Device, fields *[]string) (*datastruct.Device, error) {
	var deviceResult *datastruct.Device
	result := DB.Table("Device").Limit(1).Where(device).Select(*fields).Find(&deviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return deviceResult, nil
		} else {
			return nil, result.Error
		}
	}
	return deviceResult, nil
}

func (v *deviceQuery) CreateDevice(device *datastruct.Device) (*datastruct.Device, error) {
	result := DB.Table("Device").Create(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func (v *deviceQuery) UpdateDevice(where *datastruct.Device, device *datastruct.Device) (*datastruct.Device, error) {
	result := DB.Table("Device").Where(where).Save(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}
