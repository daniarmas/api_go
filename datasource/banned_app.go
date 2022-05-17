package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedAppDatasource interface {
	GetBannedApp(tx *gorm.DB, where *models.BannedApp) (*models.BannedApp, error)
}

type bannedAppDatasource struct{}

func (v *bannedAppDatasource) GetBannedApp(tx *gorm.DB, where *models.BannedApp) (*models.BannedApp, error) {
	var res *models.BannedApp
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("banned App not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
