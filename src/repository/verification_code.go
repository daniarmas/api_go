package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type VerificationCodeQuery interface {
	GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*[]datastruct.VerificationCode, error)
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	DeleteVerificationCode(verificationCode *datastruct.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	result := DB.Table("VerificationCode").Create(&verificationCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *verificationCodeQuery) GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*[]datastruct.VerificationCode, error) {
	var verificationCodeResult *[]datastruct.VerificationCode
	result := DB.Table("VerificationCode").Limit(1).Where(verificationCode).Select(*fields).Find(&verificationCodeResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return verificationCodeResult, nil
		} else {
			return nil, result.Error
		}
	}
	return verificationCodeResult, nil
}

func (v *verificationCodeQuery) DeleteVerificationCode(verificationCode *datastruct.VerificationCode) error {
	var verificationCodeResult *[]datastruct.VerificationCode
	result := DB.Table("VerificationCode").Where(verificationCode).Delete(&verificationCodeResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
