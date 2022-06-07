package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedDeviceDatasource interface {
	GetBannedDevice(tx *gorm.DB, where *models.BannedDevice, fields *[]string) (*models.BannedDevice, error)
}

type bannedDeviceDatasource struct{}

func (v *bannedDeviceDatasource) GetBannedDevice(tx *gorm.DB, where *models.BannedDevice, fields *[]string) (*models.BannedDevice, error) {
	var res *models.BannedDevice
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
