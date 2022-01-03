package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"gorm.io/gorm"
)

type BannedUserQuery interface {
	GetBannedUser(tx *gorm.DB, bannedUser *datastruct.BannedUser, fields *[]string) (*datastruct.BannedUser, error)
	// ListItem() ([]datastruct.Item, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type bannedUserQuery struct{}

func (i *bannedUserQuery) GetBannedUser(tx *gorm.DB, bannedUser *datastruct.BannedUser, fields *[]string) (*datastruct.BannedUser, error) {
	var bannedUserResult *datastruct.BannedUser
	result := tx.Table("BannedUser").Limit(1).Where(bannedUser).Select(*fields).Find(&bannedUserResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return bannedUserResult, nil
		} else {
			return nil, result.Error
		}
	}
	return bannedUserResult, nil
}
