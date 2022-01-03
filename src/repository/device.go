package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"gorm.io/gorm"
)

type DeviceQuery interface {
	GetDevice(tx *gorm.DB, device *datastruct.Device, fields *[]string) (*datastruct.Device, error)
	// ListItem() ([]datastruct.Item, error)
	CreateDevice(tx *gorm.DB, device *datastruct.Device) (*datastruct.Device, error)
	UpdateDevice(tx *gorm.DB, where *datastruct.Device, device *datastruct.Device) (*datastruct.Device, error)
	// DeleteItem(id int64) error
}

type deviceQuery struct{}

func (i *deviceQuery) GetDevice(tx *gorm.DB, device *datastruct.Device, fields *[]string) (*datastruct.Device, error) {
	var deviceResult *datastruct.Device
	result := tx.Table("Device").Limit(1).Where(device).Select(*fields).Find(&deviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return deviceResult, nil
		} else {
			return nil, result.Error
		}
	}
	return deviceResult, nil
}

func (v *deviceQuery) CreateDevice(tx *gorm.DB, device *datastruct.Device) (*datastruct.Device, error) {
	result := tx.Table("Device").Create(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func (v *deviceQuery) UpdateDevice(tx *gorm.DB, where *datastruct.Device, device *datastruct.Device) (*datastruct.Device, error) {
	result := tx.Table("Device").Where(where).Save(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}
