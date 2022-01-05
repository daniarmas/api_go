package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type RefreshTokenQuery interface {
	GetRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*datastruct.RefreshToken, error)
	// ListItem() ([]datastruct.Item, error)
	CreateRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error)
	// UpdateDevice(device *datastruct.Device) error
	DeleteRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) error
}

type refreshTokenQuery struct{}

func (v *refreshTokenQuery) CreateRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error) {
	result := tx.Table("RefreshToken").Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshToken, nil
}

func (r *refreshTokenQuery) DeleteRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) error {
	var refreshTokenResult *[]datastruct.RefreshToken
	result := tx.Table("RefreshToken").Where(refreshToken).Delete(&refreshTokenResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *refreshTokenQuery) GetRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*datastruct.RefreshToken, error) {
	var refreshTokenResult *datastruct.RefreshToken
	var result *gorm.DB
	if fields != nil {
		result = tx.Table("RefreshToken").Limit(1).Where(refreshToken).Select(*fields).Find(&refreshTokenResult)
	} else {
		result = tx.Table("RefreshToken").Limit(1).Where(refreshToken).Find(&refreshTokenResult)
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return refreshTokenResult, nil
		} else {
			return nil, result.Error
		}
	}
	return refreshTokenResult, nil
}
