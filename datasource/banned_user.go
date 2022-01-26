package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserDatasource interface {
	GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error)
}

type bannedUserDatasource struct{}

func (i *bannedUserDatasource) GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error) {
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
