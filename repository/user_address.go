package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserAddressRepository interface {
	ListUserAddress(tx *gorm.DB, where *models.UserAddress, fields *[]string) (*[]models.UserAddress, error)
	CreateUserAddress(tx *gorm.DB, data *models.UserAddress) (*models.UserAddress, error)
	UpdateUserAddress(tx *gorm.DB, where *models.UserAddress, data *models.UserAddress) (*models.UserAddress, error)
	GetUserAddress(tx *gorm.DB, where *models.UserAddress) (*models.UserAddress, error)
}

type userAddressRepository struct{}

func (i *userAddressRepository) ListUserAddress(tx *gorm.DB, where *models.UserAddress, fields *[]string) (*[]models.UserAddress, error) {
	result, err := Datasource.NewUserAddressDatasource().ListUserAddress(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *userAddressRepository) GetUserAddress(tx *gorm.DB, where *models.UserAddress) (*models.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().GetUserAddress(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) CreateUserAddress(tx *gorm.DB, data *models.UserAddress) (*models.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().CreateUserAddress(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) UserAddress(tx *gorm.DB, data *models.UserAddress) (*models.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().CreateUserAddress(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) UpdateUserAddress(tx *gorm.DB, where *models.UserAddress, data *models.UserAddress) (*models.UserAddress, error) {
	result, err := Datasource.NewUserAddressDatasource().UpdateUserAddress(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
