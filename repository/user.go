package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserQuery interface {
	GetUser(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error)
	GetUserWithAddress(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error)
	CreateUser(tx *gorm.DB, data *models.User) (*models.User, error)
	UpdateUser(tx *gorm.DB, where *models.User, data *models.User) (*models.User, error)
}

type userQuery struct{}

func (u *userQuery) GetUserWithAddress(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	res, err := Datasource.NewUserDatasource().GetUserWithAddress(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userQuery) GetUser(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	res, err := Datasource.NewUserDatasource().GetUser(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userQuery) CreateUser(tx *gorm.DB, data *models.User) (*models.User, error) {
	res, err := Datasource.NewUserDatasource().CreateUser(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userQuery) UpdateUser(tx *gorm.DB, where *models.User, data *models.User) (*models.User, error) {
	res, err := Datasource.NewUserDatasource().UpdateUser(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
