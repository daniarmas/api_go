package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type BannedUserQuery interface {
	GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error)
}

type bannedUserQuery struct{}

func (i *bannedUserQuery) GetBannedUser(tx *gorm.DB, bannedUser *models.BannedUser, fields *[]string) (*models.BannedUser, error) {
	result, err := Datasource.NewBannedUserDatasource().GetBannedUser(tx, bannedUser, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
