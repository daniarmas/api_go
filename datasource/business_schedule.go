package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessScheduleDatasource interface {
	GetBusinessSchedule(tx *gorm.DB, where *models.BusinessSchedule) (*models.BusinessSchedule, error)
}

type businessScheduleDatasource struct{}

func (v *businessScheduleDatasource) GetBusinessSchedule(tx *gorm.DB, where *models.BusinessSchedule) (*models.BusinessSchedule, error) {
	var response *models.BusinessSchedule
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
