package repository

import (
	"errors"

	"github.com/daniarmas/api_go/src/datastruct"
)

type VerificationCodeQuery interface {
	GetVerificationCode(code string, email string, verificationCodeType string, deviceId string) error
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	var user []datastruct.User
	var bannedUser []datastruct.BannedUser
	var bannedDevice []datastruct.BannedDevice
	var verificationCodeI datastruct.VerificationCode
	userResult := DB.Table("User").Limit(1).Where("email = ?", verificationCode.Email).Find(&user)
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
	bannedUserResult := DB.Table("BannedUser").Limit(1).Where("email = ?", verificationCode.Email).Find(&bannedUser)
	if bannedUserResult.Error != nil {
		return bannedUserResult.Error
	}
	if len(bannedUser) != 0 {
		return errors.New("banned user")
	}
	bannedDeviceResult := DB.Table("BannedDevice").Limit(1).Where("device_id = ?", verificationCode.DeviceId).Find(&bannedDevice)
	if bannedDeviceResult.Error != nil {
		return bannedDeviceResult.Error
	}
	if len(bannedDevice) != 0 {
		return errors.New("banned device")
	}
	DB.Table("VerificationCode").Where(&datastruct.VerificationCode{Email: verificationCode.Email, Type: verificationCode.Type, DeviceId: verificationCode.DeviceId}).Delete(&verificationCodeI)
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
