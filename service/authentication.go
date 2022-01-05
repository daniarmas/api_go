package service

import (
	"errors"

	"github.com/daniarmas/api_go/datastruct"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	GetVerificationCode(verificationCode *datastruct.VerificationCode, fields *[]string) (*datastruct.VerificationCode, error)
	SignIn(verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
	SignUp(fullname *string, alias *string, verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
	UserExists(email *string) error
	CheckSession(metadata *metadata.MD) (*[]string, error)
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
			if *user == (datastruct.User{}) {
				return errors.New("user not found")
			}
		case "SignUp":
			if *user != (datastruct.User{}) {
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

func (v *authenticationService) SignUp(fullname *string, alias *string, verificationCode *datastruct.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error) {
	var verificationCodeRes *datastruct.VerificationCode
	var userRes *datastruct.User
	var bannedUserRes *datastruct.BannedUser
	var bannedDeviceRes *datastruct.BannedDevice
	var deviceRes *datastruct.Device
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr, createUserErr error
	var refreshTokenRes *datastruct.RefreshToken
	var authorizationTokenRes *datastruct.AuthorizationToken
	var createUserRes *datastruct.User
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &datastruct.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceId: verificationCode.DeviceId, Type: "SignUp"}, &[]string{"id"})
		if verificationCodeErr != nil {
			return verificationCodeErr
		} else if *verificationCodeRes == (datastruct.VerificationCode{}) {
			return errors.New("verification code not found")
		}
		userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &datastruct.User{Email: verificationCode.Email}, nil)
		if userErr != nil {
			return userErr
		} else if *userRes != (datastruct.User{}) {
			return errors.New("user exists")
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
		createUserRes, createUserErr = v.dao.NewUserQuery().CreateUser(tx, &datastruct.User{Email: verificationCode.Email, Alias: *alias, IsLegalAge: true, FullName: *fullname})
		if createUserErr != nil {
			return createUserErr
		}
		deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &datastruct.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &datastruct.RefreshToken{UserFk: createUserRes.ID, DeviceFk: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &datastruct.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: createUserRes.ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
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
	return &dto.SignIn{AuthorizationToken: *jwtAuthorizationTokenRes, RefreshToken: *jwtRefreshTokenRes, User: *createUserRes}, nil
}

func (v *authenticationService) UserExists(email *string) error {
	var userRes *datastruct.User
	var userErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &datastruct.User{Email: *email}, &[]string{"id"})
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if err != nil {
		return err
	} else if *userRes != (datastruct.User{}) {
		return errors.New("user already exists")
	}
	return nil
}

func (v *authenticationService) CheckSession(metadata *metadata.MD) (*[]string, error) {
	var userRes *datastruct.User
	var bannedUserRes *datastruct.BannedUser
	var bannedDeviceRes *datastruct.BannedDevice
	var deviceRes *datastruct.Device
	var bannedDeviceErr, deviceErr, userErr, bannedUserErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		if len(metadata.Get("authorization")) != 0 && metadata.Get("authorization")[0] != "" {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if *deviceRes == (datastruct.Device{}) {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if *deviceRes != (datastruct.Device{}) {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0]}, &datastruct.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &datastruct.BannedDevice{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if bannedDeviceErr != nil {
				return bannedDeviceErr
			} else if *bannedDeviceRes != (datastruct.BannedDevice{}) {
				return errors.New("device banned")
			}
			authorizationTokenParseRes, authorizationTokenParseErr := v.dao.NewTokenQuery().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
			if authorizationTokenParseErr != nil {
				switch authorizationTokenParseErr.Error() {
				case "Token is expired":
					return errors.New("authorizationtoken expired")
				case "signature is invalid":
					return errors.New("signature is invalid")
				case "token contains an invalid number of segments":
					return errors.New("token contains an invalid number of segments")
				default:
					return authorizationTokenParseErr
				}
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &datastruct.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk", "refresh_token_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if *authorizationTokenRes == (datastruct.AuthorizationToken{}) {
				return errors.New("unauthenticated")
			}
			refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &datastruct.RefreshToken{ID: authorizationTokenRes.RefreshTokenFk}, &[]string{"id"})
			if refreshTokenErr != nil {
				return refreshTokenErr
			} else if *refreshTokenRes == (datastruct.RefreshToken{}) {
				return errors.New("unauthenticated")
			}
			userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &datastruct.User{ID: authorizationTokenRes.UserFk}, &[]string{"id"})
			if userErr != nil {
				return userErr
			} else if *userRes == (datastruct.User{}) {
				return errors.New("user not found")
			}
			bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &datastruct.BannedUser{UserFk: authorizationTokenRes.UserFk}, &[]string{"id"})
			if bannedUserErr != nil {
				return bannedUserErr
			} else if *bannedUserRes != (datastruct.BannedUser{}) {
				return errors.New("user banned")
			}
			return nil
		} else {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if *deviceRes == (datastruct.Device{}) {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if *deviceRes != (datastruct.Device{}) {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &datastruct.Device{DeviceId: metadata.Get("deviceid")[0]}, &datastruct.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &datastruct.BannedDevice{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if bannedDeviceErr != nil {
				return bannedDeviceErr
			} else if *bannedDeviceRes != (datastruct.BannedDevice{}) {
				return errors.New("device banned")
			}
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return &[]string{}, nil
}
