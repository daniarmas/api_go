package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApplicationDatasource interface {
	CreateApplication(tx *gorm.DB, data *entity.Application) (*entity.Application, error)
	GetApplication(tx *gorm.DB, where *entity.Application) (*entity.Application, error)
	ListApplication(tx *gorm.DB, where *entity.Application, cursor *time.Time) (*[]entity.Application, error)
	DeleteApplication(tx *gorm.DB, where *entity.Application, ids *[]uuid.UUID) (*[]entity.Application, error)
}

type applicationDatasource struct{}

func (i *applicationDatasource) DeleteApplication(tx *gorm.DB, where *entity.Application, ids *[]uuid.UUID) (*[]entity.Application, error) {
	var res *[]entity.Application
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

func (i *applicationDatasource) ListApplication(tx *gorm.DB, where *entity.Application, cursor *time.Time) (*[]entity.Application, error) {
	var res []entity.Application
	result := tx.Model(&entity.Application{}).Limit(11).Where(where).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (v *applicationDatasource) CreateApplication(tx *gorm.DB, data *entity.Application) (*entity.Application, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *applicationDatasource) GetApplication(tx *gorm.DB, where *entity.Application) (*entity.Application, error) {
	var res *entity.Application
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
