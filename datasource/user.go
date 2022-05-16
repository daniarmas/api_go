package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDatasource interface {
	GetUser(tx *gorm.DB, user *models.User) (*models.User, error)
	GetUserWithPermission(tx *gorm.DB, user *models.User) (*models.User, error)
	GetUserWithAddress(tx *gorm.DB, user *models.User, fields *[]string) (*models.User, error)
	CreateUser(tx *gorm.DB, user *models.User) (*models.User, error)
	UpdateUser(tx *gorm.DB, where *models.User, data *models.User) (*models.User, error)
}

type userDatasource struct{}

func (u *userDatasource) GetUserWithAddress(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var userResult *models.User
	var userAddressResult []models.UserAddress
	var userAddressErr error
	var result *gorm.DB
	query := fmt.Sprintf("SELECT id, tag, residence_type, building_number, house_number, description, user_id, province_id, municipality_id, create_time, update_time, ST_AsEWKB(coordinates) AS coordinates FROM user_address WHERE user_address.user_id = '%v';", where.ID)
	userAddressErr = tx.Raw(query).Scan(&userAddressResult).Error
	if userAddressErr != nil {
		return nil, userAddressErr
	}
	result = tx.Joins("Company").Where(where).Take(&userResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return userResult, nil
		} else {
			return nil, result.Error
		}
	}
	userResult.UserAddress = userAddressResult
	return userResult, nil
}

func (u *userDatasource) GetUser(tx *gorm.DB, where *models.User) (*models.User, error) {
	var userResult *models.User
	result := tx.Where(where).Take(&userResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return userResult, nil
}

func (u *userDatasource) GetUserWithPermission(tx *gorm.DB, where *models.User) (*models.User, error) {
	var userResult *models.User
	result := tx.Preload("Permission").Where(where).Take(&userResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return userResult, nil
}

func (u *userDatasource) CreateUser(tx *gorm.DB, user *models.User) (*models.User, error) {
	result := tx.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (v *userDatasource) UpdateUser(tx *gorm.DB, where *models.User, data *models.User) (*models.User, error) {
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
