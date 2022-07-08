package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BannedUserRepository interface {
	GetBannedUser(tx *gorm.DB, where *entity.BannedUser, fields *[]string) (*entity.BannedUser, error)
}

type bannedUserRepository struct{}

func (i *bannedUserRepository) GetBannedUser(tx *gorm.DB, where *entity.BannedUser, fields *[]string) (*entity.BannedUser, error) {
	res, err := Datasource.NewBannedUserDatasource().GetBannedUser(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
