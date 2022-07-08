package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnionBusinessRoleAndPermissionDatasource interface {
	ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission) (*[]entity.UnionBusinessRoleAndPermission, error)
	ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.UnionBusinessRoleAndPermission, error)
	CreateUnionBusinessRoleAndPermission(tx *gorm.DB, where *[]entity.UnionBusinessRoleAndPermission) (*[]entity.UnionBusinessRoleAndPermission, error)
	DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndPermission, error)
	DeleteUnionBusinessRoleAndPermissionByPermissionIds(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndPermission, error)
}

type unionBusinessRoleAndPermissionDatasource struct{}

func (i *unionBusinessRoleAndPermissionDatasource) ListUnionBusinessRoleAndPermissionInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.UnionBusinessRoleAndPermission, error) {
	var UnionBusinessRoleAndPermissions []entity.UnionBusinessRoleAndPermission
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

func (i *unionBusinessRoleAndPermissionDatasource) ListUnionBusinessRoleAndPermission(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission) (*[]entity.UnionBusinessRoleAndPermission, error) {
	var res []entity.UnionBusinessRoleAndPermission
	result := tx.Model(&entity.UnionBusinessRoleAndPermission{}).Where(where).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *unionBusinessRoleAndPermissionDatasource) CreateUnionBusinessRoleAndPermission(tx *gorm.DB, data *[]entity.UnionBusinessRoleAndPermission) (*[]entity.UnionBusinessRoleAndPermission, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *unionBusinessRoleAndPermissionDatasource) DeleteUnionBusinessRoleAndPermission(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndPermission, error) {
	var res *[]entity.UnionBusinessRoleAndPermission
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

func (v *unionBusinessRoleAndPermissionDatasource) DeleteUnionBusinessRoleAndPermissionByPermissionIds(tx *gorm.DB, where *entity.UnionBusinessRoleAndPermission, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndPermission, error) {
	var res *[]entity.UnionBusinessRoleAndPermission
	result := tx.Clauses(clause.Returning{}).Where(where).Where(`permission_id IN ?`, *ids).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}
