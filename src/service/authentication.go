package service

import (
	"errors"

	"github.com/daniarmas/api_go/src/datastruct"
	"github.com/daniarmas/api_go/src/repository"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*[]datastruct.VerificationCode, error)
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	user, _ := v.dao.NewUserQuery().GetUser(&datastruct.User{Email: verificationCode.Email})
	switch verificationCode.Type {
	case "SignIn", "ChangeUserEmail":
		if len(*user) == 0 {
			return errors.New("user not found")
		}
	case "SignUp":
		if len(*user) != 0 {
			return errors.New("user already exists")
		}
	}
	bannedUserResult, bannedUserError := v.dao.NewBannedUserQuery().GetBannedUser(&datastruct.BannedUser{Email: verificationCode.Email})
	if bannedUserError != nil {
		return bannedUserError
	}
	if len(*bannedUserResult) != 0 {
		return errors.New("banned user")
	}
	bannedDeviceResult, bannedDeviceError := v.dao.NewBannedDeviceQuery().GetBannedDevice(&datastruct.BannedDevice{DeviceId: verificationCode.DeviceId})
	if bannedDeviceError != nil {
		return bannedDeviceError
	}
	if len(*bannedDeviceResult) != 0 {
		return errors.New("banned device")
	}
	v.dao.NewVerificationCodeQuery().DeleteVerificationCode(&datastruct.VerificationCode{Email: verificationCode.Email, Type: verificationCode.Type, DeviceId: verificationCode.DeviceId})
	verificationCodeResult := v.dao.NewVerificationCodeQuery().CreateVerificationCode(verificationCode)
	if verificationCodeResult != nil {
		return verificationCodeResult
	}
	return nil
}

func (v *authenticationService) GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*[]datastruct.VerificationCode, error) {
	var response *[]datastruct.VerificationCode
	result, err := v.dao.NewVerificationCodeQuery().GetVerificationCode(verificationCode, fields)
	if err != nil {
		return response, err
	}
	return result, nil
}
