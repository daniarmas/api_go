package repository

import (
	"github.com/daniarmas/api_go/datastruct"
	"gorm.io/gorm"
)

type VerificationCodeQuery interface {
	GetVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode, fields *[]string) (*datastruct.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode) error
	DeleteVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode) error {
	result := tx.Create(&verificationCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *verificationCodeQuery) GetVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode, fields *[]string) (*datastruct.VerificationCode, error) {
	var verificationCodeResult *datastruct.VerificationCode
	result := tx.Limit(1).Where(verificationCode).Select(*fields).Find(&verificationCodeResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return verificationCodeResult, nil
		} else {
			return nil, result.Error
		}
	}
	return verificationCodeResult, nil
}

func (v *verificationCodeQuery) DeleteVerificationCode(tx *gorm.DB, verificationCode *datastruct.VerificationCode) error {
	var verificationCodeResult *[]datastruct.VerificationCode
	result := tx.Where(verificationCode).Delete(&verificationCodeResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
