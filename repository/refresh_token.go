package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenQuery interface {
	GetRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*datastruct.RefreshToken, error)
	// ListRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*[]datastruct.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error)
	// UpdateDevice(device *datastruct.Device) error
	DeleteRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*[]datastruct.RefreshToken, error)
}

type refreshTokenQuery struct{}

func (v *refreshTokenQuery) CreateRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken) (*datastruct.RefreshToken, error) {
	result := tx.Create(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshToken, nil
}

func (r *refreshTokenQuery) DeleteRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*[]datastruct.RefreshToken, error) {
	var refreshTokenResultSlice *[]datastruct.RefreshToken
	result := tx.Clauses(clause.Returning{}).Where(refreshToken).Delete(&refreshTokenResultSlice)
	if result.Error != nil {
		return nil, result.Error
	}
	return refreshTokenResultSlice, nil
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

// func (r *refreshTokenQuery) ListRefreshToken(tx *gorm.DB, refreshToken *datastruct.RefreshToken, fields *[]string) (*[]datastruct.RefreshToken, error) {
// 	var refreshTokenResult *[]datastruct.RefreshToken
// 	var result *gorm.DB
// 	if fields != nil {
// 		result = tx.Table("RefreshToken").Where(refreshToken).Select(*fields).Find(&refreshTokenResult)
// 	} else {
// 		result = tx.Table("RefreshToken").Where(refreshToken).Find(&refreshTokenResult)
// 	}
// 	if result.Error != nil {
// 		if result.Error.Error() == "record not found" {
// 			return refreshTokenResult, nil
// 		} else {
// 			return nil, result.Error
// 		}
// 	}
// 	return refreshTokenResult, nil
// }
