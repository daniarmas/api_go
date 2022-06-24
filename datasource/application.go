package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApplicationDatasource interface {
	CreateApplication(tx *gorm.DB, data *models.Application) (*models.Application, error)
	GetApplication(tx *gorm.DB, where *models.Application, fields *[]string) (*models.Application, error)
	ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error)
	DeleteApplication(tx *gorm.DB, where *models.Application, ids *[]uuid.UUID) (*[]models.Application, error)
}

type applicationDatasource struct{}

func (i *applicationDatasource) DeleteApplication(tx *gorm.DB, where *models.Application, ids *[]uuid.UUID) (*[]models.Application, error) {
	var res *[]models.Application
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (i *applicationDatasource) ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error) {
	var res []models.Application
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Model(&models.Application{}).Select(*selectFields).Limit(11).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
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
