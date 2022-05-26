package datasource

import (
	"errors"
	"fmt"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDatasource interface {
	GetUser(tx *gorm.DB, user *models.User, fields *[]string) (*models.User, error)
	GetUserWithAddress(tx *gorm.DB, user *models.User, fields *[]string) (*models.User, error)
	CreateUser(tx *gorm.DB, data *models.User) (*models.User, error)
	UpdateUser(tx *gorm.DB, where *models.User, data *models.User) (*models.User, error)
}

type userDatasource struct{}

func (u *userDatasource) GetUserWithAddress(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var res *models.User
	var userAddressResult []models.UserAddress
	var userAddressErr error
	var result *gorm.DB
	query := fmt.Sprintf("SELECT id, tag, residence_type, address, instructions, number, user_id, province_id, municipality_id, create_time, update_time, ST_AsEWKB(coordinates) AS coordinates FROM user_address WHERE user_address.user_id = '%v';", where.ID)
	userAddressErr = tx.Raw(query).Scan(&userAddressResult).Error
	if userAddressErr != nil {
		return nil, userAddressErr
	}
	result = tx.Preload("UserPermissions").Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return res, nil
		} else {
			return nil, result.Error
		}
	}
	res.UserAddress = userAddressResult
	return res, nil
}

func (u *userDatasource) GetUser(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var res *models.User
	result := tx.Where(where).Select(*fields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (u *userDatasource) CreateUser(tx *gorm.DB, data *models.User) (*models.User, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
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
