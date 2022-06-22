package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ApplicationDatasource interface {
	CreateApplication(tx *gorm.DB, data *models.Application) (*models.Application, error)
	GetApplication(tx *gorm.DB, where *models.Application, fields *[]string) (*models.Application, error)
	ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error)
}

type applicationDatasource struct{}

func (i *applicationDatasource) ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error) {
	var res []models.Application
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Select(*selectFields).Limit(11).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (v *applicationDatasource) CreateApplication(tx *gorm.DB, data *models.Application) (*models.Application, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *applicationDatasource) GetApplication(tx *gorm.DB, where *models.Application, fields *[]string) (*models.Application, error) {
	var res *models.Application
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
