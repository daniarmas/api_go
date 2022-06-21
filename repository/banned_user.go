package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserRepository interface {
	GetBannedUser(tx *gorm.DB, where *models.BannedUser, fields *[]string) (*models.BannedUser, error)
}

type bannedUserRepository struct{}

func (i *bannedUserRepository) GetBannedUser(tx *gorm.DB, where *models.BannedUser, fields *[]string) (*models.BannedUser, error) {
	res, err := Datasource.NewBannedUserDatasource().GetBannedUser(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}
