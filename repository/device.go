package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeviceQuery interface {
	GetDevice(tx *gorm.DB, device *models.Device, fields *[]string) (*models.Device, error)
	// ListItem() ([]models.Item, error)
	CreateDevice(tx *gorm.DB, device *models.Device) (*models.Device, error)
	UpdateDevice(tx *gorm.DB, where *models.Device, device *models.Device) (*models.Device, error)
	// DeleteItem(id int64) error
}

type deviceQuery struct{}

func (i *deviceQuery) GetDevice(tx *gorm.DB, device *models.Device, fields *[]string) (*models.Device, error) {
	var deviceResult *models.Device
	result := tx.Limit(1).Where(device).Select(*fields).Find(&deviceResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return deviceResult, nil
		} else {
			return nil, result.Error
		}
	}
	return deviceResult, nil
}

func (v *deviceQuery) CreateDevice(tx *gorm.DB, device *models.Device) (*models.Device, error) {
	result := tx.Create(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func (v *deviceQuery) UpdateDevice(tx *gorm.DB, where *models.Device, device *models.Device) (*models.Device, error) {
	result := tx.Where(where).Updates(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}
