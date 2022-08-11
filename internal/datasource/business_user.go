package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessUserDatasource interface {
	GetBusinessUser(tx *gorm.DB, where *entity.BusinessUser) (*entity.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error)
}

type businessUserDatasource struct{}

func (v *businessUserDatasource) CreateBusinessUser(tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *businessUserDatasource) GetBusinessUser(tx *gorm.DB, where *entity.BusinessUser) (*entity.BusinessUser, error) {
	var res *entity.BusinessUser
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

func (v *businessUserDatasource) DeleteBusinessUser(tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error) {
	var res *[]entity.BusinessUser
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
