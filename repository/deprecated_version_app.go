package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeprecatedVersionAppRepository interface {
	GetDeprecatedVersionApp(tx *gorm.DB, where *models.DeprecatedVersionApp) (*models.DeprecatedVersionApp, error)
}

type deprecatedVersionAppRepository struct{}

func (i *deprecatedVersionAppRepository) GetDeprecatedVersionApp(tx *gorm.DB, where *models.DeprecatedVersionApp) (*models.DeprecatedVersionApp, error) {
	res, err := Datasource.NewDeprecatedVersionAppDatasource().GetDeprecatedVersionApp(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
