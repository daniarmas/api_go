package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BannedDeviceDatasource interface {
	GetBannedDevice(tx *gorm.DB, where *entity.BannedDevice, fields *[]string) (*entity.BannedDevice, error)
}

type bannedDeviceDatasource struct{}

func (v *bannedDeviceDatasource) GetBannedDevice(tx *gorm.DB, where *entity.BannedDevice, fields *[]string) (*entity.BannedDevice, error) {
	var res *entity.BannedDevice
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
