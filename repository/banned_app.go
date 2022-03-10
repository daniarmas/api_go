package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedAppRepository interface {
	GetBannedApp(tx *gorm.DB, where *models.BannedApp) (*models.BannedApp, error)
}

type bannedAppRepository struct{}

func (i *bannedAppRepository) GetBannedApp(tx *gorm.DB, where *models.BannedApp) (*models.BannedApp, error) {
	result, err := Datasource.NewBannedAppDatasource().GetBannedApp(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
