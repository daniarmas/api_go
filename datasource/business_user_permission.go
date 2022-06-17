package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessUserPermissionDatasource interface {
	GetBusinessUserPermission(tx *gorm.DB, where *models.BusinessUserPermission, fields *[]string) (*models.BusinessUserPermission, error)
}

type businessUserPermissionDatasource struct{}

func (i *businessUserPermissionDatasource) GetBusinessUserPermission(tx *gorm.DB, where *models.BusinessUserPermission, fields *[]string) (*models.BusinessUserPermission, error) {
	var res *models.BusinessUserPermission
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
