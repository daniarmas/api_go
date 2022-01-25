package repository

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserQuery interface {
	GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error)
	// ListItem() ([]models.Item, error)
	// CreateItem(answer models.Item) (*int64, error)
	// UpdateItem(answer models.Item) (*models.Item, error)
	// DeleteItem(id int64) error
}

type bannedUserQuery struct{}

func (i *bannedUserQuery) GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error) {
	var bannedUserResult *models.BannedUser
	result := tx.Where(bannedUser).Select(*fields).Take(&bannedUserResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return bannedUserResult, nil
}
