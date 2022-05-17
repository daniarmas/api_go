package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserQuery interface {
	GetBannedUser(tx *gorm.DB, where *models.BannedUser) (*models.BannedUser, error)
}

type bannedUserQuery struct{}

func (i *bannedUserQuery) GetBannedUser(tx *gorm.DB, where *models.BannedUser) (*models.BannedUser, error) {
	res, err := Datasource.NewBannedUserDatasource().GetBannedUser(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}
