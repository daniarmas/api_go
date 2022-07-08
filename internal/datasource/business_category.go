package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessCategoryDatasource interface {
	GetBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, fields *[]string) (*entity.BusinessCategory, error)
	CreateBusinessCategory(tx *gorm.DB, data *entity.BusinessCategory) (*entity.BusinessCategory, error)
	DeleteBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, ids *[]uuid.UUID) (*[]entity.BusinessCategory, error)
	ListBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, fields *[]string) (*[]entity.BusinessCategory, error)
}

type businessCategoryDatasource struct{}

func (i *businessCategoryDatasource) ListBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, fields *[]string) (*[]entity.BusinessCategory, error) {
	var res []entity.BusinessCategory
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (v *businessCategoryDatasource) CreateBusinessCategory(tx *gorm.DB, data *entity.BusinessCategory) (*entity.BusinessCategory, error) {
	var existBusinessCategory *entity.BusinessCategory
	existResult := tx.Where("name = ?", data.Name).Select("id").Take(&existBusinessCategory)
	if existResult.Error != nil && existResult.Error.Error() != "record not found" {
		return nil, existResult.Error
	}
	if existResult.Error.Error() == "record not found" {
		result := tx.Create(&data)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		return nil, errors.New("record exists")
	}
	return data, nil
}

func (v *businessCategoryDatasource) GetBusinessCategory(tx *gorm.DB, BusinessCategory *entity.BusinessCategory, fields *[]string) (*entity.BusinessCategory, error) {
	var res *entity.BusinessCategory
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(BusinessCategory).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *businessCategoryDatasource) DeleteBusinessCategory(tx *gorm.DB, where *entity.BusinessCategory, ids *[]uuid.UUID) (*[]entity.BusinessCategory, error) {
	var res *[]entity.BusinessCategory
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
