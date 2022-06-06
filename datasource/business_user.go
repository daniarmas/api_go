package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessUserDatasource interface {
	GetBusinessUser(tx *gorm.DB, where *models.BusinessUser, fields *[]string) (*models.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *models.BusinessUser) (*models.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser, ids *[]uuid.UUID) (*[]models.BusinessUser, error)
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
	var res *models.BusinessUser
	selectFields := &[]string{"*"}
	if fields == nil {
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

func (v *businessUserDatasource) DeleteBusinessUser(tx *gorm.DB, where *models.BusinessUser, ids *[]uuid.UUID) (*[]models.BusinessUser, error) {
	var res *[]models.BusinessUser
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
