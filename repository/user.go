package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type UserQuery interface {
	GetUser(tx *gorm.DB, user *models.User, fields *[]string) (*models.User, error)
	CreateUser(tx *gorm.DB, user *models.User) (*models.User, error)
}

type userQuery struct{}

func (u *userQuery) GetUser(tx *gorm.DB, where *models.User, fields *[]string) (*models.User, error) {
	var userResult *models.User
	var result *gorm.DB
	if fields != nil {
		result = tx.Limit(1).Where(where).Select(*fields).Find(&userResult)
	} else {
		result = tx.Limit(1).Where(where).Find(&userResult)
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return userResult, nil
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
