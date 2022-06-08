package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessCategoryDatasource interface {
	GetBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, fields *[]string) (*models.BusinessCategory, error)
	CreateBusinessCategory(tx *gorm.DB, data *models.BusinessCategory) (*models.BusinessCategory, error)
	DeleteBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, ids *[]uuid.UUID) (*[]models.BusinessCategory, error)
}

type businessCategoryDatasource struct{}

func (v *businessCategoryDatasource) CreateBusinessCategory(tx *gorm.DB, data *models.BusinessCategory) (*models.BusinessCategory, error) {
	var existBusinessCategory *models.BusinessCategory
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

func (v *businessCategoryDatasource) GetBusinessCategory(tx *gorm.DB, BusinessCategory *models.BusinessCategory, fields *[]string) (*models.BusinessCategory, error) {
	var res *models.BusinessCategory
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

func (v *businessCategoryDatasource) DeleteBusinessCategory(tx *gorm.DB, where *models.BusinessCategory, ids *[]uuid.UUID) (*[]models.BusinessCategory, error) {
	var res *[]models.BusinessCategory
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
