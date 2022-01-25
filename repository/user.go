package repository

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserQuery interface {
	GetUser(tx *gorm.DB, user *models.User) (*models.User, error)
	GetUserWithAddress(tx *gorm.DB, user *models.User, fields *[]string) (*models.User, error)
	CreateUser(tx *gorm.DB, user *models.User) (*models.User, error)
}

type userQuery struct{}

func (u *userQuery) GetUserWithAddress(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var userResult *models.User
	var userAddressResult *[]models.UserAddress
	var userAddressErr error
	var result *gorm.DB
	query := fmt.Sprintf("SELECT id, tag, residence_type, building_number, house_number, description, user_fk, province_fk, municipality_fk, create_time, update_time, ST_AsEWKB(coordinates) AS coordinates FROM user_address WHERE user_address.user_fk = '%v';", where.ID)
	userAddressErr = tx.Raw(query).Scan(&userAddressResult).Error
	if userAddressErr != nil {
		return nil, userAddressErr
	}
	result = tx.Where(where).Take(&userResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return userResult, nil
		} else {
			return nil, result.Error
		}
	}
	userResult.UserAddress = *userAddressResult
	return userResult, nil
}

func (u *userQuery) GetUser(tx *gorm.DB, where *models.User) (*models.User, error) {
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

func (u *userQuery) CreateUser(tx *gorm.DB, user *models.User) (*models.User, error) {
	result := tx.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
