package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAddressRepository interface {
	ListUserAddress(tx *gorm.DB, where *entity.UserAddress, fields *[]string) (*[]entity.UserAddress, error)
	CreateUserAddress(tx *gorm.DB, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddress(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddressByUserId(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddressSelected(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	GetUserAddress(tx *gorm.DB, where *entity.UserAddress) (*entity.UserAddress, error)
	DeleteUserAddress(tx *gorm.DB, where *entity.UserAddress, ids *[]uuid.UUID) (*[]entity.UserAddress, error)
}

type userAddressRepository struct{}

func (i *userAddressRepository) ListUserAddress(tx *gorm.DB, where *entity.UserAddress, fields *[]string) (*[]entity.UserAddress, error) {
	result, err := Datasource.NewUserAddressDatasource().ListUserAddress(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *userAddressRepository) DeleteUserAddress(tx *gorm.DB, where *entity.UserAddress, ids *[]uuid.UUID) (*[]entity.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().DeleteUserAddress(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) GetUserAddress(tx *gorm.DB, where *entity.UserAddress) (*entity.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().GetUserAddress(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) CreateUserAddress(tx *gorm.DB, data *entity.UserAddress) (*entity.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().CreateUserAddress(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) UpdateUserAddressByUserId(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().UpdateUserAddressByUserId(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) UpdateUserAddressSelected(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	res, err := Datasource.NewUserAddressDatasource().UpdateUserAddressSelected(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *userAddressRepository) UpdateUserAddress(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	result, err := Datasource.NewUserAddressDatasource().UpdateUserAddress(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
