package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"gorm.io/gorm"
)

type UserQuery interface {
	GetUser(tx *gorm.DB, user *datastruct.User, fields *[]string) (*datastruct.User, error)
	// ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type userQuery struct{}

func (i *userQuery) GetUser(tx *gorm.DB, user *datastruct.User, fields *[]string) (*datastruct.User, error) {
	var userResult *datastruct.User
	var result *gorm.DB
	if fields != nil {
		result = tx.Table("User").Limit(1).Where(user).Select(*fields).Find(&userResult)
	} else {
		result = tx.Table("User").Limit(1).Where(user).Find(&userResult)
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
