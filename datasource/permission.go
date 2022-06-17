package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessPermissionDatasource interface {
	GetBusinessPermission(tx *gorm.DB, where *models.BusinessPermission, fields *[]string) (*models.BusinessPermission, error)
}

type businessPermissionDatasource struct{}

func (i *businessPermissionDatasource) GetBusinessPermission(tx *gorm.DB, where *models.BusinessPermission, fields *[]string) (*models.BusinessPermission, error) {
	var res *models.BusinessPermission
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
