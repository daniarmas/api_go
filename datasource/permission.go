package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type PermissionDatasource interface {
	PermissionExists(tx *gorm.DB, where *models.Permission) error
}

type permissionDatasource struct{}

func (i *permissionDatasource) PermissionExists(tx *gorm.DB, where *models.Permission) error {
	var permissionResult *models.Permission
	result := tx.Where(where).Select("id").Take(&permissionResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return errors.New("record not found")
		} else {
			return result.Error
		}
	}
	return nil
}
