package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc/metadata"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error)
	GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error)
	SignIn(verificationCode *models.VerificationCode, metadata *metadata.MD) (*dto.SignIn, error)
	SignUp(fullname *string, alias *string, verificationCode *models.VerificationCode, signUpType *string, metadata *metadata.MD) (*dto.SignIn, error)
	SignOut(all *bool, authorizationTokenid *string, metadata *metadata.MD) error
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

func (v *authenticationService) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error) {
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		user, err := v.dao.NewUserQuery().GetUser(tx, &models.User{Email: req.Email})
		if err != nil {
			if err.Error() == "record not found" && (req.Type.String() == "SignIn" || req.Type.String() == "ChangeUserEmail") {
				return errors.New("user not found")
			} else if user != nil && req.Type.String() == "SignUp" {
				return errors.New("user already exists")
			}
		}
		bannedUserResult, err := v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: req.Email})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if bannedUserResult != nil {
			return errors.New("banned user")
		}
		bannedDeviceResult, err := v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: *meta.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if bannedDeviceResult != nil {
			return errors.New("banned device")
		}
		// bannedAppRes, bannedAppErr := v.dao.NewBannedAppRepository().GetBannedApp(tx, &models.BannedApp{Version: metadata.Get("appversion")[0]})
		// if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
		// 	return bannedAppErr
		// } else if bannedAppRes != nil {
		// 	return errors.New("app banned")
		// }
		v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *meta.DeviceIdentifier})
		verificationCodeResult := v.dao.NewVerificationCodeQuery().CreateVerificationCode(tx, &models.VerificationCode{Code: utils.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceIdentifier: *meta.DeviceIdentifier, CreateTime: time.Now(), UpdateTime: time.Now()})
		if verificationCodeResult != nil {
			return verificationCodeResult
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
}

func (v *authenticationService) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error) {
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Code: req.Code, Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
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
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceIdentifier: verificationCode.DeviceIdentifier, Type: "SignIn"}, &[]string{"id"})
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
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: verificationCode.Email})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: verificationCode.DeviceIdentifier})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		bannedAppRes, bannedAppErr := v.dao.NewBannedAppRepository().GetBannedApp(tx, &models.BannedApp{Version: metadata.Get("appversion")[0]})
		if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
			return bannedAppErr
		} else if bannedAppRes != nil {
			return errors.New("app banned")
		}
		deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceIdentifier: verificationCode.DeviceIdentifier})
		if deleteVerificationCodeErr != nil {
			return deleteVerificationCodeErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier}, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		if len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &(*deleteRefreshTokenRes)[0].ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &refreshTokenRes.ID, UserId: &userRes.ID, DeviceId: &deviceRes.ID, App: &metadata.Get("app")[0], AppVersion: &metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{TokenId: &refreshTokenRes.ID}
		jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
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
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Code: verificationCode.Code, DeviceIdentifier: verificationCode.DeviceIdentifier, Type: "SignUp"}, &[]string{"id"})
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
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: verificationCode.Email})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: verificationCode.ID})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		bannedAppRes, bannedAppErr := v.dao.NewBannedAppRepository().GetBannedApp(tx, &models.BannedApp{Version: metadata.Get("appversion")[0]})
		if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
			return bannedAppErr
		} else if bannedAppRes != nil {
			return errors.New("app banned")
		}
		deleteVerificationCodeErr := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: verificationCode.Email, Type: "SignIn", DeviceIdentifier: verificationCode.DeviceIdentifier})
		if deleteVerificationCodeErr != nil {
			return deleteVerificationCodeErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier}, &models.Device{DeviceIdentifier: verificationCode.DeviceIdentifier, Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		createUserRes, createUserErr = v.dao.NewUserQuery().CreateUser(tx, &models.User{Email: verificationCode.Email, IsLegalAge: true, FullName: *fullname})
		if createUserErr != nil {
			return createUserErr
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		if len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &(*deleteRefreshTokenRes)[0].ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: createUserRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &refreshTokenRes.ID, UserId: &createUserRes.ID, DeviceId: &deviceRes.ID, App: &metadata.Get("app")[0], AppVersion: &metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{TokenId: &refreshTokenRes.ID}
		jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		var isBusinessOwner bool = false
		if verificationCodeRes.Type == "SignUpBusinessOwner" {
			isBusinessOwner = true
		}
		if *signUpType == "SignUpBusiness" {
			_, createBusinessUserErr := v.dao.NewBusinessUserRepository().CreateBusinessUser(tx, &models.BusinessUser{IsBusinessOwner: isBusinessOwner, UserId: createUserRes.ID})
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

func (v *authenticationService) CheckSession(metadata *metadata.MD) (*[]string, error) {
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var bannedDeviceErr, deviceErr, userErr, bannedUserErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if len(metadata.Get("authorization")) != 0 && metadata.Get("authorization")[0] != "" {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if deviceRes == nil {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if deviceRes != nil {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: metadata.Get("deviceid")[0]})
			if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
				return bannedDeviceErr
			} else if bannedDeviceRes != nil {
				return errors.New("device banned")
			}
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
			authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: *authorizationTokenRes.RefreshTokenId})
			if refreshTokenErr != nil {
				return refreshTokenErr
			} else if refreshTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserId})
			if userErr != nil {
				return userErr
			} else if userRes == nil {
				return errors.New("user not found")
			}
			bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{UserId: *authorizationTokenRes.UserId})
			if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
				return bannedUserErr
			} else if bannedUserRes != nil {
				return errors.New("user banned")
			}
			bannedAppRes, bannedAppErr := v.dao.NewBannedAppRepository().GetBannedApp(tx, &models.BannedApp{Version: metadata.Get("appversion")[0]})
			if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
				return bannedAppErr
			} else if bannedAppRes != nil {
				return errors.New("app banned")
			}
			return nil
		} else {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &[]string{"id"})
			if deviceErr != nil {
				return deviceErr
			} else if deviceRes == nil {
				deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			} else if deviceRes != nil {
				_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
				if deviceErr != nil {
					return deviceErr
				}
			}
			bannedAppRes, bannedAppErr := v.dao.NewBannedAppRepository().GetBannedApp(tx, &models.BannedApp{Version: metadata.Get("appversion")[0]})
			if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
				return bannedAppErr
			} else if bannedAppRes != nil {
				return errors.New("app banned")
			}
			bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: metadata.Get("deviceid")[0]})
			if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
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

func (v *authenticationService) SignOut(all *bool, authorizationTokenid *string, metadata *metadata.MD) error {
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if *all {
			authorizationWebToken := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
			authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(authorizationWebToken)
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: authorizationWebToken.TokenId})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			var refreshTokenIds []string
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: *authorizationTokenRes.UserId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			for _, e := range *deleteRefreshTokenRes {
				refreshTokenIds = append(refreshTokenIds, e.ID.String())
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, "refresh_token_id IN ?", refreshTokenIds)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else if *authorizationTokenid != "" {
			jwtTokenAuthorization := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
			authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtTokenAuthorization)
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
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtTokenAuthorization.TokenId})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtTokenAuthorization.TokenId})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			} else if authorizationTokenRes.UserId != authorizationTokenByReqRes.UserId {
				return errors.New("permission denied")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: *authorizationTokenByReqRes.UserId, DeviceId: *authorizationTokenByReqRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &(*deleteRefreshTokenRes)[0].ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			jwtTokenAuthorization := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
			authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtTokenAuthorization)
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtTokenAuthorization.TokenId})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: *authorizationTokenRes.UserId, DeviceId: *authorizationTokenRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &(*deleteRefreshTokenRes)[0].ID}, nil)
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
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
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
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listSessionRes, listSessionErr = v.dao.NewSessionQuery().ListSession(tx, &models.Session{UserId: *authorizationTokenRes.UserId}, nil)
		if listSessionErr != nil {
			return listSessionErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.ListSessionResponse{Sessions: listSessionRes, ActualDeviceId: *authorizationTokenRes.DeviceId}, nil
}

func (v *authenticationService) RefreshToken(refreshToken *string, metadata *metadata.MD) (*dto.RefreshToken, error) {
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	var jwtAuthorizationTokenErr, jwtRefreshTokenErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{Token: refreshToken}
		refreshTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtRefreshToken(jwtRefreshToken)
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
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: *jwtRefreshToken.TokenId})
		if refreshTokenErr != nil {
			return refreshTokenErr
		} else if refreshTokenRes == nil {
			return errors.New("unauthenticated")
		}
		userRes, userErr := v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: refreshTokenRes.UserId}, &[]string{"id"})
		if userErr != nil {
			return userErr
		} else if userRes == nil {
			return errors.New("user not found")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{ID: refreshTokenRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &(*deleteRefreshTokenRes)[0].ID}, nil)
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		deviceRes, deviceErr := v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if *deviceRes == (models.Device{}) {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0], Platform: metadata.Get("platform")[0], SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0], Model: metadata.Get("model")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		} else if *deviceRes != (models.Device{}) {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: metadata.Get("deviceid")[0]}, &models.Device{SystemVersion: metadata.Get("systemversion")[0], FirebaseCloudMessagingId: metadata.Get("firebasecloudmessagingid")[0]})
			if deviceErr != nil {
				return deviceErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenId: &refreshTokenRes.ID, UserId: &userRes.ID, DeviceId: &deviceRes.ID, App: &metadata.Get("app")[0], AppVersion: &metadata.Get("appversion")[0]})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshTokenNew := &datasource.JsonWebTokenMetadata{TokenId: &refreshTokenRes.ID}
		jwtAuthorizationTokenNew := &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshTokenNew)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationTokenNew)
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
