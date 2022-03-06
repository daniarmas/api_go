package usecase

import (
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *models.VerificationCode) error
	GetVerificationCode(verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	SignIn(verificationCode *models.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
	SignUp(fullname *string, alias *string, verificationCode *models.VerificationCode, signUpType *string, metadata *metadata.MD) (*dto.SignIn, error)
	SignOut(all *bool, authorizationTokenFk *string, metadata *metadata.MD) error
	UserExists(alias *string) error
	CheckSession(metadata *metadata.MD) (*[]string, error)
	ListSession(metadata *metadata.MD) (*dto.ListSessionResponse, error)
	RefreshToken(refreshToken *string, metadata *metadata.MD) (*dto.RefreshToken, error)
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) CreateVerificationCode(verificationCode *models.VerificationCode) error {
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		user, userErr := v.dao.NewUserQuery().GetUser(tx, &models.User{Email: verificationCode.Email})
		if userErr != nil {
			if userErr.Error() == "record not found" && (verificationCode.Type == "SignIn" || verificationCode.Type == "ChangeUserEmail") {
				return errors.New("user not found")
			} else if user != nil && verificationCode.Type == "SignUp" {
				return errors.New("user already exists")
			}
		}
		bannedUserResult, bannedUserError := v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
		if bannedUserError != nil && bannedUserError.Error() != "record not found" {
			return bannedUserError
		}
		if bannedUserResult != nil {
			return errors.New("banned user")
		}
		bannedDeviceResult, bannedDeviceError := v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if bannedDeviceError != nil && bannedUserError.Error() != "record not found" {
			return bannedDeviceError
		}
		if bannedDeviceResult != nil {
			return errors.New("banned device")
		}
		v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Type: verificationCode.Type, DeviceId: verificationCode.DeviceId})
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

func (v *authenticationService) GetVerificationCode(verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	txErr := datasource.DB.Transaction(func(tx *gorm.DB) error {
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

func (v *authenticationService) SignIn(verificationCode *models.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error) {
	var verificationCodeRes *models.VerificationCode
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr error
	var refreshTokenRes *models.RefreshToken
	var authorizationTokenRes *models.AuthorizationToken
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceId: verificationCode.DeviceId, Type: "SignIn"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeRes == nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserQuery().GetUserWithPermission(tx, &models.User{Email: verificationCode.Email})
		if userErr != nil {
			switch userErr.Error() {
			case "record not found":
				return errors.New("user not found")
			default:
				return userErr
			}
		}
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceId: verificationCode.DeviceId})
		if deleteVerificationCodeErr != nil {
			return deleteVerificationCodeErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId}, &models.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID}, &[]string{"id"})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		if len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: userRes.ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(&authorizationTokenId)
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

func (v *authenticationService) SignUp(fullname *string, alias *string, verificationCode *models.VerificationCode, signUpType *string, metadata *metadata.MD) (*dto.SignIn, error) {
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var verificationCodeRes *models.VerificationCode
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr, createUserErr error
	var refreshTokenRes *models.RefreshToken
	var authorizationTokenRes *models.AuthorizationToken
	var createUserRes *models.User
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceId: verificationCode.DeviceId, Type: "SignUp"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeErr != nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{Email: verificationCode.Email}, nil)
		if userErr != nil {
			return userErr
		} else if userRes.Email != "" {
			return errors.New("user exists")
		}
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: verificationCode.Email}, &[]string{"id"})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceId: verificationCode.DeviceId})
		if deleteVerificationCodeErr != nil {
			return deleteVerificationCodeErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceId: verificationCode.DeviceId}, &models.Device{DeviceId: verificationCode.DeviceId, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		createUserRes, createUserErr = v.dao.NewUserQuery().CreateUser(tx, &models.User{Email: verificationCode.Email, Alias: *alias, IsLegalAge: true, FullName: *fullname})
		if createUserErr != nil {
			return createUserErr
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID}, &[]string{"id"})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		if len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserFk: createUserRes.ID, DeviceFk: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: createUserRes.ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(&authorizationTokenId)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		var isBusinessOwner bool = false
		if verificationCodeRes.Type == "SignUpBusinessOwner" {
			isBusinessOwner = true
		}
		if *signUpType == "SignUpBusiness" {
			_, createBusinessUserErr := v.dao.NewBusinessUserRepository().CreateBusinessUser(tx, &models.BusinessUser{IsBusinessOwner: isBusinessOwner, UserFk: createUserRes.ID})
			if createBusinessUserErr != nil {
				return createBusinessUserErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.SignIn{AuthorizationToken: *jwtAuthorizationTokenRes, RefreshToken: *jwtRefreshTokenRes, User: *createUserRes}, nil
}

func (v *authenticationService) UserExists(alias *string) error {
	var userRes *models.User
	var userErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		userRes, userErr = v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{Alias: *alias}, &[]string{"id"})
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if err != nil {
		return err
	} else if userRes != nil {
		return errors.New("user already exists")
	}
	return nil
}

func (v *authenticationService) CheckSession(metadata *metadata.MD) (*[]string, error) {
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var bannedDeviceErr, deviceErr, userErr, bannedUserErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if len(metadata.Get("authorization")) != 0 && metadata.Get("authorization")[0] != "" {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if *deviceRes == (models.Device{}) {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if *deviceRes != (models.Device{}) {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if bannedDeviceErr != nil {
				return bannedDeviceErr
			} else if bannedDeviceRes != nil {
				return errors.New("device banned")
			}
			authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk", "refresh_token_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: authorizationTokenRes.RefreshTokenFk}, &[]string{"id"})
			if refreshTokenErr != nil {
				return refreshTokenErr
			} else if refreshTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr = v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: authorizationTokenRes.UserFk}, &[]string{"id"})
			if userErr != nil {
				return userErr
			} else if userRes == nil {
				return errors.New("user not found")
			}
			bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{UserFk: authorizationTokenRes.UserFk}, &[]string{"id"})
			if bannedUserErr != nil {
				return bannedUserErr
			} else if bannedUserRes != nil {
				return errors.New("user banned")
			}
			return nil
		} else {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if *deviceRes == (models.Device{}) {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if *deviceRes != (models.Device{}) {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if bannedDeviceErr != nil {
				return bannedDeviceErr
			} else if bannedDeviceRes != nil {
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

func (v *authenticationService) SignOut(all *bool, authorizationTokenFk *string, metadata *metadata.MD) error {
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if *all {
			authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			var refreshTokenIds []string
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenRes.UserFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			for _, e := range *deleteRefreshTokenRes {
				refreshTokenIds = append(refreshTokenIds, e.ID.String())
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationTokenIn(tx, "refresh_token_fk IN ?", &refreshTokenIds)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else if *authorizationTokenFk != "" {
			authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
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
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenFk)}, &[]string{"id", "user_fk", "device_fk"})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			} else if authorizationTokenRes.UserFk != authorizationTokenByReqRes.UserFk {
				return errors.New("permission denied")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenByReqRes.UserFk, DeviceFk: authorizationTokenByReqRes.DeviceFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk", "device_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenRes.UserFk, DeviceFk: authorizationTokenRes.DeviceFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *authenticationService) ListSession(metadata *metadata.MD) (*dto.ListSessionResponse, error) {
	var listSessionRes *[]models.Session
	var authorizationTokenRes *models.AuthorizationToken
	var authorizationTokenErr error
	var listSessionErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
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
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk", "device_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listSessionRes, listSessionErr = v.dao.NewSessionQuery().ListSession(tx, &models.Session{UserFk: authorizationTokenRes.UserFk}, nil)
		if listSessionErr != nil {
			return listSessionErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.ListSessionResponse{Sessions: listSessionRes, ActualDeviceId: authorizationTokenRes.DeviceFk}, nil
}

func (v *authenticationService) RefreshToken(refreshToken *string, metadata *metadata.MD) (*dto.RefreshToken, error) {
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	var jwtAuthorizationTokenErr, jwtRefreshTokenErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		refreshTokenParseRes, refreshTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtRefreshToken(refreshToken)
		if refreshTokenParseErr != nil {
			switch refreshTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("refreshtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return refreshTokenParseErr
			}
		}
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: uuid.MustParse(*refreshTokenParseRes)}, &[]string{"id", "user_fk", "device_fk"})
		if refreshTokenErr != nil {
			return refreshTokenErr
		} else if refreshTokenRes == nil {
			return errors.New("unauthenticated")
		}
		userRes, userErr := v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: refreshTokenRes.UserFk}, &[]string{"id"})
		if userErr != nil {
			return userErr
		} else if userRes == nil {
			return errors.New("user not found")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{ID: refreshTokenRes.ID}, &[]string{"id"})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		deviceRes, deviceErr := v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if *deviceRes == (models.Device{}) {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else if *deviceRes != (models.Device{}) {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceId: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID, DeviceFk: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: userRes.ID, DeviceFk: deviceRes.ID, App: metadata.Get("app")[0], AppVersion: metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(&authorizationTokenId)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.RefreshToken{RefreshToken: *jwtRefreshTokenRes, AuthorizationToken: *jwtAuthorizationTokenRes}, nil
}
