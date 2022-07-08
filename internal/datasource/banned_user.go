package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BannedUserDatasource interface {
	GetBannedUser(tx *gorm.DB, where *entity.BannedUser, fields *[]string) (*entity.BannedUser, error)
}

type bannedUserDatasource struct{}

func (v *bannedUserDatasource) GetBannedUser(tx *gorm.DB, where *entity.BannedUser, fields *[]string) (*entity.BannedUser, error) {
	var res *entity.BannedUser
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
