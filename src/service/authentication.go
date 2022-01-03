package service

import (
	"errors"

	"github.com/daniarmas/api_go/src/datastruct"
	"github.com/daniarmas/api_go/src/dto"
	"github.com/daniarmas/api_go/src/repository"
	"google.golang.org/grpc/metadata"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*[]datastruct.VerificationCode, error)
	SignIn(verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	user, _ := v.dao.NewUserQuery().GetUser(&datastruct.User{Email: verificationCode.Email}, &[]string{"id"})
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
	bannedUserResult, bannedUserError := v.dao.NewBannedUserQuery().GetBannedUser(&datastruct.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
	if bannedUserError != nil {
		return bannedUserError
	}
	if len(*bannedUserResult) != 0 {
		return errors.New("banned user")
	}
	bannedDeviceResult, bannedDeviceError := v.dao.NewBannedDeviceQuery().GetBannedDevice(&datastruct.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
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

func (v *authenticationService) SignIn(verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error) {
	verificationCodeRes, verificationCodeErr := v.dao.NewVerificationCodeQuery().GetVerificationCode(&datastruct.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceId: verificationCode.DeviceId, Type: "SignIn"}, &[]string{"id"})
	if verificationCodeErr != nil {
		return nil, verificationCodeErr
	} else if len(*verificationCodeRes) == 0 {
		return nil, errors.New("verification code not found")
	}
	userRes, userErr := v.dao.NewUserQuery().GetUser(&datastruct.User{Email: verificationCode.Email}, nil)
	if userErr != nil {
		return nil, userErr
	} else if len(*userRes) == 0 {
		return nil, errors.New("user not found")
	}
	bannedUserRes, bannedUserErr := v.dao.NewBannedUserQuery().GetBannedUser(&datastruct.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
	if bannedUserErr != nil {
		return nil, bannedUserErr
	} else if len(*bannedUserRes) != 0 {
		return nil, errors.New("user banned")
	}
	bannedDeviceRes, bannedDeviceErr := v.dao.NewBannedDeviceQuery().GetBannedDevice(&datastruct.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
	if bannedDeviceErr != nil {
		return nil, bannedDeviceErr
	} else if len(*bannedDeviceRes) != 0 {
		return nil, errors.New("device banned")
	}
	deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(&datastruct.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceId: verificationCode.DeviceId})
	if deleteVerificationCodeErr != nil {
		return nil, deleteVerificationCodeErr
	}
	deviceRes, deviceErr := v.dao.NewDeviceQuery().GetDevice(&datastruct.Device{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
	if deviceErr != nil {
		return nil, deviceErr
	} else if deviceRes == nil {
		deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(&datastruct.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
		if deviceErr != nil {
			return nil, deviceErr
		}
	} else if deviceRes != nil {
		_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(&datastruct.Device{DeviceId: verificationCode.DeviceId}, &datastruct.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
		if deviceErr != nil {
			return nil, deviceErr
		}
	}
	deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(&datastruct.RefreshToken{UserFk: (*userRes)[0].ID, DeviceFk: deviceRes.ID})
	if deleteRefreshTokenErr != nil {
		return nil, deleteRefreshTokenErr
	}
	refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().CreateRefreshToken(&datastruct.RefreshToken{UserFk: (*userRes)[0].ID, DeviceFk: deviceRes.ID})
	if refreshTokenErr != nil {
		return nil, refreshTokenErr
	}
	authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(&datastruct.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: (*userRes)[0].ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
	if authorizationTokenErr != nil {
		return nil, authorizationTokenErr
	}
	authorizationTokenId := authorizationTokenRes.ID.String()
	refreshTokenId := refreshTokenRes.ID.String()
	jwtRefreshTokenRes, jwtRefreshTokenErr := v.dao.NewTokenQuery().CreateJwtRefreshToken(&refreshTokenId)
	if jwtRefreshTokenErr != nil {
		return nil, jwtRefreshTokenErr
	}
	jwtAuthorizationTokenRes, jwtAuthorizationTokenErr := v.dao.NewTokenQuery().CreateJwtAuthorizationToken(&authorizationTokenId)
	if jwtAuthorizationTokenErr != nil {
		return nil, jwtAuthorizationTokenErr
	}
	return &dto.SignIn{AuthorizationToken: *jwtAuthorizationTokenRes, RefreshToken: *jwtRefreshTokenRes, User: (*userRes)[0]}, nil
}
