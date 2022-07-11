package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserConfigurationDatasource interface {
	ListUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error)
	CreateUserConfiguration(tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	UpdateUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	GetUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error)
	DeleteUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error)
}

type userConfigurationDatasource struct{}

func (i *userConfigurationDatasource) ListUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error) {
	var res []entity.UserConfiguration
	result := tx.Where(where).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *userConfigurationDatasource) CreateUserConfiguration(tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *userConfigurationDatasource) UpdateUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}

func (i *userConfigurationDatasource) GetUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	var res *entity.UserConfiguration
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

func (v *userConfigurationDatasource) DeleteUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error) {
	var res *[]entity.UserConfiguration
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}
