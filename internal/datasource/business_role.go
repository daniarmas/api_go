package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessRoleDatasource interface {
	CreateBusinessRole(tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error)
	UpdateBusinessRole(tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error)
	GetBusinessRole(tx *gorm.DB, where *entity.BusinessRole, fields *[]string) (*entity.BusinessRole, error)
	ListBusinessRole(tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time, fields *[]string) (*[]entity.BusinessRole, error)
	DeleteBusinessRole(tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error)
}

type businessRoleDatasource struct{}

func (v *businessRoleDatasource) UpdateBusinessRole(tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return data, nil
}

func (v *businessRoleDatasource) CreateBusinessRole(tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *businessRoleDatasource) DeleteBusinessRole(tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error) {
	var res *[]entity.BusinessRole
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

func (i *businessRoleDatasource) ListBusinessRole(tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time, fields *[]string) (*[]entity.BusinessRole, error) {
	var res []entity.BusinessRole
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Model(&entity.BusinessRole{}).Select(*selectFields).Limit(11).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *businessRoleDatasource) GetBusinessRole(tx *gorm.DB, where *entity.BusinessRole, fields *[]string) (*entity.BusinessRole, error) {
	var res *entity.BusinessRole
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
