package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessScheduleRepository interface {
	GetBusinessSchedule(tx *gorm.DB, where *entity.BusinessSchedule) (*entity.BusinessSchedule, error)
	BusinessIsOpen(tx *gorm.DB, where *entity.BusinessSchedule, orderType string) (bool, error)
}

type businessScheduleRepository struct{}

func (i *businessScheduleRepository) GetBusinessSchedule(tx *gorm.DB, where *entity.BusinessSchedule) (*entity.BusinessSchedule, error) {
	result, err := Datasource.NewBusinessScheduleDatasource().GetBusinessSchedule(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessScheduleRepository) BusinessIsOpen(tx *gorm.DB, where *entity.BusinessSchedule, orderType string) (bool, error) {
	result, err := Datasource.NewBusinessScheduleDatasource().BusinessIsOpen(tx, where, orderType)
	if err != nil {
		return false, err
	}
	return result, nil
}
