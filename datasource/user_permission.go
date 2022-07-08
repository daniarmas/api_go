package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPermissionDatasource interface {
	CreateUserPermission(tx *gorm.DB, data *[]entity.UserPermission) (*[]entity.UserPermission, error)
	GetUserPermission(tx *gorm.DB, where *entity.UserPermission, fields *[]string) (*entity.UserPermission, error)
	DeleteUserPermission(tx *gorm.DB, where *entity.UserPermission, ids *[]uuid.UUID) (*[]entity.UserPermission, error)
	DeleteUserPermissionByPermissionId(tx *gorm.DB, permissionIds *[]uuid.UUID) (*[]entity.UserPermission, error)
	DeleteUserPermissionByBusinessRoleId(tx *gorm.DB, where *entity.UserPermission) (*[]entity.UserPermission, error)
}

type userPermissionDatasource struct{}

func (i *userPermissionDatasource) DeleteUserPermissionByBusinessRoleId(tx *gorm.DB, where *entity.UserPermission) (*[]entity.UserPermission, error) {
	var res *[]entity.UserPermission
	result := tx.Clauses(clause.Returning{}).Where(`business_role_id = ?`, where.BusinessRoleId).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (i *userPermissionDatasource) CreateUserPermission(tx *gorm.DB, data *[]entity.UserPermission) (*[]entity.UserPermission, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *userPermissionDatasource) GetUserPermission(tx *gorm.DB, where *entity.UserPermission, fields *[]string) (*entity.UserPermission, error) {
	var res *entity.UserPermission
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

func (i *userPermissionDatasource) DeleteUserPermissionByPermissionId(tx *gorm.DB, permissionIds *[]uuid.UUID) (*[]entity.UserPermission, error) {
	var res *[]entity.UserPermission
	result := tx.Clauses(clause.Returning{}).Where(`permission_id IN ?`, *permissionIds).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *userPermissionDatasource) DeleteUserPermission(tx *gorm.DB, where *entity.UserPermission, ids *[]uuid.UUID) (*[]entity.UserPermission, error) {
	var res *[]entity.UserPermission
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
