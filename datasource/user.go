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
	result := tx.Preload("UserPermissions").Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return res, nil
		} else {
			return nil, result.Error
		}
	}
	query := fmt.Sprintf("SELECT id, tag, address, instructions, number, user_id, province_id, municipality_id, create_time, update_time, ST_AsEWKB(coordinates) AS coordinates FROM user_address WHERE user_address.user_id = '%v';", res.ID)
	userAddressErr = tx.Raw(query).Scan(&userAddressResult).Error
	if userAddressErr != nil {
		return nil, userAddressErr
	}
	res.UserAddress = userAddressResult
	return res, nil
}

func (u *userDatasource) GetUser(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var res *models.User
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
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
