package datasource

import (
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UnionBusinessRoleAndUserDatasource interface {
	ListUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, cursor *time.Time, fields *[]string) (*[]entity.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error)
	ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.UnionBusinessRoleAndUser, error)
	CreateUnionBusinessRoleAndUser(tx *gorm.DB, where *[]entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error)
	DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error)
	GetUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*entity.UnionBusinessRoleAndUser, error)
}

type unionBusinessRoleAndUserDatasource struct{}

func (i *unionBusinessRoleAndUserDatasource) GetUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*entity.UnionBusinessRoleAndUser, error) {
	var res *entity.UnionBusinessRoleAndUser
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

func (i *unionBusinessRoleAndUserDatasource) ListUnionBusinessRoleAndUserInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]entity.UnionBusinessRoleAndUser, error) {
	var UnionBusinessRoleAndUsers []entity.UnionBusinessRoleAndUser
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where("id IN ?", ids).Select(*selectFields).Find(&UnionBusinessRoleAndUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	return &UnionBusinessRoleAndUsers, nil
}

func (i *unionBusinessRoleAndUserDatasource) ListUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, cursor *time.Time, fields *[]string) (*[]entity.UnionBusinessRoleAndUser, error) {
	var res []entity.UnionBusinessRoleAndUser
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Model(&entity.UnionBusinessRoleAndUser{}).Select(*selectFields).Limit(11).Where("create_time < ?", cursor).Order("create_time desc").Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *unionBusinessRoleAndUserDatasource) ListUnionBusinessRoleAndUserAll(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error) {
	var res []entity.UnionBusinessRoleAndUser
	result := tx.Model(&entity.UnionBusinessRoleAndUser{}).Where(where).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *unionBusinessRoleAndUserDatasource) CreateUnionBusinessRoleAndUser(tx *gorm.DB, data *[]entity.UnionBusinessRoleAndUser) (*[]entity.UnionBusinessRoleAndUser, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *unionBusinessRoleAndUserDatasource) DeleteUnionBusinessRoleAndUser(tx *gorm.DB, where *entity.UnionBusinessRoleAndUser, ids *[]uuid.UUID) (*[]entity.UnionBusinessRoleAndUser, error) {
	var res *[]entity.UnionBusinessRoleAndUser
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
