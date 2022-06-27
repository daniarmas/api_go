package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnionBusinessRoleAndPermissionDatasource interface {
	ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, cursor *time.Time, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error)
	ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error)
	CreateUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission) (*models.UnionBusinessRoleAndPermission, error)
	DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error)
}

type unionBusinessRoleAndPermissionDatasource struct{}

func (i *unionBusinessRoleAndPermissionDatasource) ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error) {
	var UnionBusinessRoleAndPermissions []models.UnionBusinessRoleAndPermission
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where("id IN ?", ids).Select(*selectFields).Find(&UnionBusinessRoleAndPermissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &UnionBusinessRoleAndPermissions, nil
}

func (i *unionBusinessRoleAndPermissionDatasource) ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, cursor *time.Time, fields *[]string) (*[]models.UnionBusinessRoleAndPermission, error) {
	var res []models.UnionBusinessRoleAndPermission
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Model(&models.UnionBusinessRoleAndPermission{}).Select(*selectFields).Limit(11).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *unionBusinessRoleAndPermissionDatasource) CreateUnionBusinessRoleAndPermission(tx *gorm.DB, data *models.UnionBusinessRoleAndPermission) (*models.UnionBusinessRoleAndPermission, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *unionBusinessRoleAndPermissionDatasource) DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *models.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]models.UnionBusinessRoleAndPermission, error) {
	var res *[]models.UnionBusinessRoleAndPermission
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
