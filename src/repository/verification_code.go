package repository

import (
	"github.com/daniarmas/api_go/src/datastruct"
)

type VerificationCodeQuery interface {
	GetVerificationCode(code string, email string, verificationCodeType string, deviceId string) error
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	var user []datastruct.User
	userResult := DB.Table("User").Where("email = ?", verificationCode.Email).Take(&user)
	switch verificationCode.Type {
	case "SignIn", "ChangeUserEmail":
		if len(user) == 0 {
			return userResult.Error
		}
	case "SignUp":
		if len(user) != 0 {
			return userResult.Error
		}
	}
	verificationCodeResult := DB.Table("VerificationCode").Create(&verificationCode)
	if verificationCodeResult.Error != nil {
		return verificationCodeResult.Error
	}
	return nil
}

func (v *verificationCodeQuery) GetVerificationCode(code string, email string, verificationCodeType string, deviceId string) error {
	var verificationCode []datastruct.VerificationCode
	verificationCodeResult := DB.Table("VerificationCode").Where("code = ?", code).Take(&verificationCode)
	if len(verificationCode) == 0 {
		return verificationCodeResult.Error
	}
	return nil
}
