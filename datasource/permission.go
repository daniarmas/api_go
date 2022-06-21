package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type PermissionDatasource interface {
	GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error)
}

type permissionDatasource struct{}

func (i *permissionDatasource) GetPermission(tx *gorm.DB, where *models.Permission, fields *[]string) (*models.Permission, error) {
	var res *models.Permission
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
