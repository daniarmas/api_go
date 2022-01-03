package service

import (
	"errors"

	"github.com/daniarmas/api_go/src/datastruct"
	"github.com/daniarmas/api_go/src/dto"
	"github.com/daniarmas/api_go/src/repository"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*datastruct.VerificationCode, error)
	SignIn(verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		user, _ := v.dao.NewUserQuery().GetUser(tx, &datastruct.User{Email: verificationCode.Email}, &[]string{"id"})
		switch verificationCode.Type {
		case "SignIn", "ChangeUserEmail":
			if user == nil {
				return errors.New("user not found")
			}
		case "SignUp":
			if user != nil {
				return errors.New("user already exists")
			}
		}
		bannedUserResult, bannedUserError := v.dao.NewBannedUserQuery().GetBannedUser(tx, &datastruct.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
		if bannedUserError != nil {
			return bannedUserError
		}
		if *bannedUserResult != (datastruct.BannedUser{}) {
			return errors.New("banned user")
		}
		bannedDeviceResult, bannedDeviceError := v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &datastruct.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if bannedDeviceError != nil {
			return bannedDeviceError
		}
		if *bannedDeviceResult != (datastruct.BannedDevice{}) {
			return errors.New("banned device")
		}
		v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &datastruct.VerificationCode{Email: verificationCode.Email, Type: verificationCode.Type, DeviceId: verificationCode.DeviceId})
		verificationCodeResult := v.dao.NewVerificationCodeQuery().CreateVerificationCode(tx, verificationCode)
		if verificationCodeResult != nil {
			return verificationCodeResult
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *authenticationService) GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*datastruct.VerificationCode, error) {
	txErr := repository.DB.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, verificationCode, fields)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return verificationCode, nil
}

func (v *authenticationService) SignIn(verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error) {
	var verificationCodeRes *datastruct.VerificationCode
	var userRes *datastruct.User
	var bannedUserRes *datastruct.BannedUser
	var bannedDeviceRes *datastruct.BannedDevice
	var deviceRes *datastruct.Device
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr error
	var refreshTokenRes *datastruct.RefreshToken
	var authorizationTokenRes *datastruct.AuthorizationToken
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &datastruct.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceId: verificationCode.DeviceId, Type: "SignIn"}, &[]string{"id"})
		if verificationCodeErr != nil {
			return verificationCodeErr
		} else if *verificationCodeRes == (datastruct.VerificationCode{}) {
			return errors.New("verification code not found")
		}
		userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &datastruct.User{Email: verificationCode.Email}, nil)
		if userErr != nil {
			return userErr
		} else if *userRes == (datastruct.User{}) {
			return errors.New("user not found")
		}
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &datastruct.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
		if bannedUserErr != nil {
			return bannedUserErr
		} else if *bannedUserRes != (datastruct.BannedUser{}) {
			return errors.New("user banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &datastruct.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if bannedDeviceErr != nil {
			return bannedDeviceErr
		} else if *bannedDeviceRes != (datastruct.BannedDevice{}) {
			return errors.New("device banned")
		}
		deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &datastruct.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceId: verificationCode.DeviceId})
		if deleteVerificationCodeErr != nil {
			return deleteVerificationCodeErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &datastruct.Device{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if *deviceRes == (datastruct.Device{}) {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &datastruct.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else if *deviceRes != (datastruct.Device{}) {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &datastruct.Device{DeviceId: verificationCode.DeviceId}, &datastruct.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &datastruct.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &datastruct.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &datastruct.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: userRes.ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = v.dao.NewTokenQuery().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = v.dao.NewTokenQuery().CreateJwtAuthorizationToken(&authorizationTokenId)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.SignIn{AuthorizationToken: *jwtAuthorizationTokenRes, RefreshToken: *jwtRefreshTokenRes, User: *userRes}, nil
}
