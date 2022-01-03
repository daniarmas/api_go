package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type RefreshTokenQuery interface {
	// GetDevice(device *datastruct.Device, fields *[]string) (*[]datastruct.Device, error)
	// ListItem() ([]datastruct.Item, error)
	CreateRefreshToken(refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error)
	// UpdateDevice(device *datastruct.Device) error
	DeleteRefreshToken(refreshToken *datastruct.RefreshToken) error
}

type refreshTokenQuery struct{}

func (v *refreshTokenQuery) CreateRefreshToken(refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error) {
	result := DB.Table("RefreshToken").Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshToken, nil
}

func (r *refreshTokenQuery) DeleteRefreshToken(refreshToken *datastruct.RefreshToken) error {
	var refreshTokenResult *[]datastruct.RefreshToken
	result := DB.Table("RefreshToken").Where(refreshToken).Delete(&refreshTokenResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
