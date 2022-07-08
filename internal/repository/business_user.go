package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessUserRepository interface {
	GetBusinessUser(tx *gorm.DB, where *entity.BusinessUser, fields *[]string) (*entity.BusinessUser, error)
	CreateBusinessUser(tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error)
	DeleteBusinessUser(tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error)
}

type businessUserRepository struct{}

func (v *businessUserRepository) CreateBusinessUser(tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().CreateBusinessUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessUserRepository) GetBusinessUser(tx *gorm.DB, where *entity.BusinessUser, fields *[]string) (*entity.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().GetBusinessUser(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *businessUserRepository) DeleteBusinessUser(tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error) {
	res, err := Datasource.NewBusinessUserDatasource().DeleteBusinessUser(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
