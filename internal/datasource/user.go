package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDatasource interface {
	GetUser(tx *gorm.DB, user *entity.User) (*entity.User, error)
	GetUserWithAddress(tx *gorm.DB, user *entity.User) (*entity.User, error)
	CreateUser(tx *gorm.DB, data *entity.User) (*entity.User, error)
	UpdateUser(tx *gorm.DB, where *entity.User, data *entity.User) (*entity.User, error)
}

type userDatasource struct{}

func (u *userDatasource) GetUserWithAddress(tx *gorm.DB, where *entity.User) (*entity.User, error) {
	var res *entity.User
	var userAddressResult []entity.UserAddress
	var userAddressErr error
	result := tx.Preload("UserPermissions").Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return res, nil
		} else {
			return nil, result.Error
		}
	}
	query := fmt.Sprintf("SELECT id, selected, name, address, instructions, number, user_id, province_id, municipality_id, create_time, update_time, ST_AsEWKB(coordinates) AS coordinates FROM user_address WHERE user_address.user_id = '%v' AND user_address.delete_time IS NULL;", res.ID)
	userAddressErr = tx.Raw(query).Scan(&userAddressResult).Error
	if userAddressErr != nil {
		return nil, userAddressErr
	}
	res.UserAddress = userAddressResult
	return res, nil
}

func (u *userDatasource) GetUser(tx *gorm.DB, where *entity.User) (*entity.User, error) {
	var res *entity.User
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

func (u *userDatasource) CreateUser(tx *gorm.DB, data *entity.User) (*entity.User, error) {
	var existUser *entity.User
	existResult := tx.Where("email = ?", data.Email).Select("id").Take(&existUser)
	if existResult.Error != nil && existResult.Error.Error() != "record not found" {
		return nil, existResult.Error
	}
	if existResult.Error.Error() == "record not found" {
		result := tx.Create(&data)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		return nil, errors.New("record exists")
	}
	return data, nil
}

func (v *userDatasource) UpdateUser(tx *gorm.DB, where *entity.User, data *entity.User) (*entity.User, error) {
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
