package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserDatasource interface {
	GetBannedUser(tx *gorm.DB, where *models.BannedUser) (*models.BannedUser, error)
}

type bannedUserDatasource struct{}

func (v *bannedUserDatasource) GetBannedUser(tx *gorm.DB, where *models.BannedUser) (*models.BannedUser, error) {
	var res *models.BannedUser
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("banned user not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
