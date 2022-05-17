package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type DeprecatedVersionAppDatasource interface {
	GetDeprecatedVersionApp(tx *gorm.DB, where *models.DeprecatedVersionApp) (*models.DeprecatedVersionApp, error)
}

type deprecatedVersionAppDatasource struct{}

func (v *deprecatedVersionAppDatasource) GetDeprecatedVersionApp(tx *gorm.DB, where *models.DeprecatedVersionApp) (*models.DeprecatedVersionApp, error) {
	var res *models.DeprecatedVersionApp
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
