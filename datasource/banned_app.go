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

func (i *bannedAppDatasource) GetBannedApp(tx *gorm.DB, where *models.BannedApp) (*models.BannedApp, error) {
	var bannedAppResult *models.BannedApp
	result := tx.Where(where).Take(&bannedAppResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return bannedAppResult, nil
}
