package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessRoleDatasource interface {
	CreateBusinessRole(tx *gorm.DB, data *models.BusinessRole) (*models.BusinessRole, error)
	GetBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*models.BusinessRole, error)
	ListBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*[]models.BusinessRole, error)
	DeleteBusinessRole(tx *gorm.DB, where *models.BusinessRole, ids *[]uuid.UUID) (*[]models.BusinessRole, error)
}

type businessRoleDatasource struct{}

func (v *businessRoleDatasource) CreateBusinessRole(tx *gorm.DB, data *models.BusinessRole) (*models.BusinessRole, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *businessRoleDatasource) DeleteBusinessRole(tx *gorm.DB, where *models.BusinessRole, ids *[]uuid.UUID) (*[]models.BusinessRole, error) {
	var res *[]models.BusinessRole
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

func (i *businessRoleDatasource) ListBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*[]models.BusinessRole, error) {
	var res []models.BusinessRole
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

func (i *businessRoleDatasource) GetBusinessRole(tx *gorm.DB, where *models.BusinessRole, fields *[]string) (*models.BusinessRole, error) {
	var res *models.BusinessRole
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
