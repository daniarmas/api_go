package repository

import (
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
	result, err := Datasource.NewUserDatasource().GetUserWithAddress(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userQuery) GetUser(tx *gorm.DB, where *models.User) (*models.User, error) {
	result, err := Datasource.NewUserDatasource().GetUser(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userQuery) CreateUser(tx *gorm.DB, user *models.User) (*models.User, error) {
	result, err := Datasource.NewUserDatasource().CreateUser(tx, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}
