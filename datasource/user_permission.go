package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserPermissionDatasource interface {
	UserPermissionExists(tx *gorm.DB, where *models.UserPermission) error
}

type userPermissionDatasource struct{}

func (i *userPermissionDatasource) UserPermissionExists(tx *gorm.DB, where *models.UserPermission) error {
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
