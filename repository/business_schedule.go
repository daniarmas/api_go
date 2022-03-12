package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BusinessScheduleRepository interface {
	GetBusinessSchedule(tx *gorm.DB, where *models.BusinessSchedule) (*models.BusinessSchedule, error)
}

type businessScheduleRepository struct{}

func (i *businessScheduleRepository) GetBusinessSchedule(tx *gorm.DB, where *models.BusinessSchedule) (*models.BusinessSchedule, error) {
	result, err := Datasource.NewBusinessScheduleDatasource().GetBusinessSchedule(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}
