package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessUserDatasource interface {
	GetBusinessUser(tx *gorm.DB, where *models.BusinessUser, fields *[]string) (*models.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser) error
}

type businessUserDatasource struct{}

func (v *businessUserDatasource) CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *businessUserDatasource) GetBusinessUser(tx *gorm.DB, where *models.BusinessUser, fields *[]string) (*models.BusinessUser, error) {
	var response *models.BusinessUser
	result := tx.Where(where).Take(&response)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return response, nil
}

func (v *businessUserDatasource) DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser) error {
	var response *[]models.BusinessUser
	result := tx.Where(where).Delete(&response)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
