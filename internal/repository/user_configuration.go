package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserConfigurationRepository interface {
	ListUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error)
	CreateUserConfiguration(tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	UpdateUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	GetUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error)
	DeleteUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error)
}

type userConfigurationRepository struct{}

func (i *userConfigurationRepository) ListUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error) {
	result, err := Datasource.NewUserConfigurationDatasource().ListUserConfiguration(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *userConfigurationRepository) DeleteUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error) {
	res, err := Datasource.NewUserConfigurationDatasource().DeleteUserConfiguration(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userConfigurationRepository) GetUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	res, err := Datasource.NewUserConfigurationDatasource().GetUserConfiguration(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userConfigurationRepository) CreateUserConfiguration(tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	res, err := Datasource.NewUserConfigurationDatasource().CreateUserConfiguration(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userConfigurationRepository) UpdateUserConfiguration(tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	result, err := Datasource.NewUserConfigurationDatasource().UpdateUserConfiguration(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
