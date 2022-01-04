package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type RefreshTokenQuery interface {
	// GetDevice(device *datastruct.Device, fields *[]string) (*[]datastruct.Device, error)
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
