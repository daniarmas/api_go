package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserPermissionDatasource interface {
	GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error)
}

type userPermissionDatasource struct{}

func (i *userPermissionDatasource) GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error) {
	var res *models.UserPermission
	result := tx.Where(where).Select(*fields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
