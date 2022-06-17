package usecase

import (
	"context"
	"fmt"

	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	smtp "github.com/daniarmas/api_go/utils/smtp"
	"github.com/google/uuid"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error)
	GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error)
	SignIn(ctx context.Context, req *pb.SignInRequest, meta *utils.ClientMetadata) (*pb.SignInResponse, error)
	SignUp(ctx context.Context, req *pb.SignUpRequest, meta *utils.ClientMetadata) (*pb.SignUpResponse, error)
	SignOut(ctx context.Context, req *pb.SignOutRequest, meta *utils.ClientMetadata) (*gp.Empty, error)
	CheckSession(ctx context.Context, meta *utils.ClientMetadata) (*[]string, error)
	ListSession(ctx context.Context, meta *utils.ClientMetadata) (*pb.ListSessionResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, meta *utils.ClientMetadata) (*pb.RefreshTokenResponse, error)
}

type authenticationService struct {
	dao    repository.DAO
	config *utils.Config
}

func NewAuthenticationService(dao repository.DAO, config *utils.Config) AuthenticationService {
	return &authenticationService{dao: dao, config: config}
}

func (v *authenticationService) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error) {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		user, err := v.dao.NewUserQuery().GetUser(tx, &models.User{Email: req.Email}, &[]string{"id"})
		if err != nil {
			if err.Error() == "record not found" && (req.Type.String() == "SignIn") {
				return errors.New("user not found")
			}
		} else if user != nil && (req.Type.String() == "SignUp" || req.Type.String() == "ChangeUserEmail") {
			return errors.New("user already exists")
		}
		bannedUserResult, err := v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: req.Email}, &[]string{"id"})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if bannedUserResult != nil {
			return errors.New("banned user")
		}
		bannedDeviceResult, err := v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if bannedDeviceResult != nil {
			return errors.New("banned device")
		}
		deprecatedVersionAppRes, err := v.dao.NewDeprecatedVersionAppRepository().GetDeprecatedVersionApp(tx, &models.DeprecatedVersionApp{Version: *meta.AppVersion}, &[]string{"id"})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if deprecatedVersionAppRes != nil {
			return errors.New("app banned")
		}
		v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *meta.DeviceIdentifier}, nil)
		createVerificationCodeRes, createVerificationCodeErr := v.dao.NewVerificationCodeQuery().CreateVerificationCode(tx, &models.VerificationCode{Code: utils.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceIdentifier: *meta.DeviceIdentifier, CreateTime: time.Now(), UpdateTime: time.Now()})
		if createVerificationCodeErr != nil {
			return createVerificationCodeErr
		}
		verificationCodeMsg := fmt.Sprintf("Su código de verificación es %s", createVerificationCodeRes.Code)
		go smtp.SendMail(req.Email, v.config.EmailAddress, v.config.EmailAddressPassword, "Código de Verificación", verificationCodeMsg, v.config)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
}

func (v *authenticationService) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, meta *utils.ClientMetadata) (*gp.Empty, error) {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
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

func (v *authenticationService) SignIn(ctx context.Context, req *pb.SignInRequest, md *utils.ClientMetadata) (*pb.SignInResponse, error) {
	var verificationCodeRes *models.VerificationCode
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr error
	var refreshTokenRes *models.RefreshToken
	var authorizationTokenRes *models.AuthorizationToken
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		bannedAppRes, bannedAppErr := v.dao.NewDeprecatedVersionAppRepository().GetDeprecatedVersionApp(tx, &models.DeprecatedVersionApp{Version: *md.AppVersion}, &[]string{})
		if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
			return bannedAppErr
		} else if bannedAppRes != nil {
			return errors.New("app banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignIn"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeRes == nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{Email: req.Email}, nil)
		if userErr != nil {
			switch userErr.Error() {
			case "record not found":
				return errors.New("user not found")
			default:
				return userErr
			}
		}
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: req.Email}, &[]string{})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		_, err := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{Email: req.Email, Type: "SignIn", DeviceIdentifier: *md.DeviceIdentifier}, nil)
		if err != nil {
			return err
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: *md.DeviceIdentifier}, &models.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID}, nil)
		if deleteRefreshTokenErr != nil && deleteRefreshTokenErr.Error() != "record not found" {
			return deleteRefreshTokenErr
		}
		if deleteRefreshTokenRes != nil && len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(ctx, tx, &models.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(ctx, tx, &models.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: userRes.ID, DeviceId: deviceRes.ID, App: md.App, AppVersion: md.AppVersion})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshToken = &datasource.JsonWebTokenMetadata{TokenId: refreshTokenRes.ID}
		jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationToken = &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	userAddress := make([]*pb.UserAddress, 0, len(userRes.UserAddress))
	permissions := make([]*pb.Permission, 0, len(userRes.BusinessUserPermissions))
	if *md.App == "Business" {
		for _, item := range userRes.BusinessUserPermissions {
			permissions = append(permissions, &pb.Permission{
				Id:         item.ID.String(),
				Name:       item.Name,
				UserId:     item.UserId.String(),
				BusinessId: item.BusinessId.String(),
				CreateTime: timestamppb.New(item.CreateTime),
				UpdateTime: timestamppb.New(item.UpdateTime),
			})
		}
	}
	for _, item := range userRes.UserAddress {
		userAddress = append(userAddress, &pb.UserAddress{
			Id:             item.ID.String(),
			Tag:            item.Tag,
			Number:         item.Number,
			Address:        item.Address,
			Instructions:   item.Instructions,
			ProvinceId:     item.ProvinceId.String(),
			MunicipalityId: item.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: item.Coordinates.Coords()[0], Longitude: item.Coordinates.Coords()[1]},
			UserId:         item.UserId.String(),
			CreateTime:     timestamppb.New(item.CreateTime),
			UpdateTime:     timestamppb.New(item.UpdateTime),
		})
	}
	go smtp.SendSignInMail(req.Email, time.Now(), v.config, md)
	var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
	if userRes.HighQualityPhoto != "" {
		highQualityPhotoUrl = v.config.UsersBulkName + "/" + userRes.HighQualityPhoto
		lowQualityPhotoUrl = v.config.UsersBulkName + "/" + userRes.LowQualityPhoto
		thumbnailUrl = v.config.UsersBulkName + "/" + userRes.Thumbnail

	}
	return &pb.SignInResponse{AuthorizationToken: *jwtAuthorizationToken.Token, RefreshToken: *jwtRefreshToken.Token, User: &pb.User{
		Id:                  userRes.ID.String(),
		FullName:            userRes.FullName,
		Email:               userRes.Email,
		HighQualityPhoto:    userRes.HighQualityPhoto,
		HighQualityPhotoUrl: highQualityPhotoUrl,
		LowQualityPhoto:     userRes.LowQualityPhoto,
		LowQualityPhotoUrl:  lowQualityPhotoUrl,
		Thumbnail:           userRes.Thumbnail,
		ThumbnailUrl:        thumbnailUrl,
		BlurHash:            userRes.BlurHash,
		Permissions:         permissions,
		UserAddress:         userAddress,
		CreateTime:          timestamppb.New(userRes.CreateTime),
		UpdateTime:          timestamppb.New(userRes.UpdateTime),
	}}, nil
}

func (v *authenticationService) SignUp(ctx context.Context, req *pb.SignUpRequest, meta *utils.ClientMetadata) (*pb.SignUpResponse, error) {
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var verificationCodeRes *models.VerificationCode
	var verificationCodeErr, userErr, bannedUserErr, bannedDeviceErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr, createUserErr error
	var refreshTokenRes *models.RefreshToken
	var authorizationTokenRes *models.AuthorizationToken
	var createUserRes *models.User
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		bannedAppRes, bannedAppErr := v.dao.NewDeprecatedVersionAppRepository().GetDeprecatedVersionApp(tx, &models.DeprecatedVersionApp{Version: *meta.AppVersion}, &[]string{})
		if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
			return bannedAppErr
		} else if bannedAppRes != nil {
			return errors.New("app banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *meta.DeviceIdentifier, Type: "SignUp"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeErr != nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &models.User{Email: req.Email}, &[]string{"id"})
		if userErr != nil && userErr.Error() != "record not found" {
			return userErr
		} else if userRes != nil {
			return errors.New("user exists")
		}
		bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{Email: req.Email}, &[]string{"id"})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		} else if bannedUserRes != nil {
			return errors.New("user banned")
		}
		_, err := v.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{ID: verificationCodeRes.ID}, nil)
		if err != nil {
			return err
		}
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier, Platform: *meta.Platform, SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId, Model: *meta.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier, Platform: *meta.Platform, SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId, Model: *meta.Model})
			if deviceErr != nil {
				return deviceErr
			}
		}
		createUserRes, createUserErr = v.dao.NewUserQuery().CreateUser(tx, &models.User{Email: req.Email, IsLegalAge: true, FullName: req.FullName})
		if createUserErr != nil {
			return createUserErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: createUserRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(ctx, tx, &models.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: createUserRes.ID, DeviceId: deviceRes.ID, App: meta.App, AppVersion: meta.AppVersion})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshToken = &datasource.JsonWebTokenMetadata{TokenId: refreshTokenRes.ID}
		jwtRefreshTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationToken = &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		jwtAuthorizationTokenErr = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		var isBusinessOwner bool = false
		if verificationCodeRes.Type == "SignUpBusinessOwner" {
			isBusinessOwner = true
		}
		if req.SignUpType.String() == "SignUpBusiness" {
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
	return &pb.SignUpResponse{AuthorizationToken: *jwtAuthorizationToken.Token, RefreshToken: *jwtRefreshToken.Token, User: &pb.User{
		Id:         createUserRes.ID.String(),
		FullName:   createUserRes.FullName,
		Email:      createUserRes.Email,
		CreateTime: timestamppb.New(createUserRes.CreateTime),
		UpdateTime: timestamppb.New(createUserRes.UpdateTime),
	}}, nil
}

func (v *authenticationService) CheckSession(ctx context.Context, meta *utils.ClientMetadata) (*[]string, error) {
	var userRes *models.User
	var bannedUserRes *models.BannedUser
	var bannedDeviceRes *models.BannedDevice
	var deviceRes *models.Device
	var bannedDeviceErr, deviceErr, userErr, bannedUserErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		deviceRes, deviceErr = v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier, Platform: *meta.Platform, SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId, Model: *meta.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &models.Device{SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId})
			if deviceErr != nil {
				return deviceErr
			}
		}
		bannedAppRes, bannedAppErr := v.dao.NewDeprecatedVersionAppRepository().GetDeprecatedVersionApp(tx, &models.DeprecatedVersionApp{Version: *meta.AppVersion}, &[]string{"id"})
		if bannedAppErr != nil && bannedAppErr.Error() != "record not found" {
			return bannedAppErr
		} else if bannedAppRes != nil {
			return errors.New("app banned")
		}
		bannedDeviceRes, bannedDeviceErr = v.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		} else if bannedDeviceRes != nil {
			return errors.New("device banned")
		}
		if meta.Authorization != nil && *meta.Authorization != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: meta.Authorization}
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: authorizationTokenRes.RefreshTokenId}, &[]string{"id"})
			if refreshTokenErr != nil {
				return refreshTokenErr
			} else if refreshTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{"id"})
			if userErr != nil {
				return userErr
			} else if userRes == nil {
				return errors.New("user not found")
			}
			bannedUserRes, bannedUserErr = v.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{UserId: authorizationTokenRes.UserId}, &[]string{"id"})
			if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
				return bannedUserErr
			} else if bannedUserRes != nil {
				return errors.New("user banned")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &[]string{}, nil
}

func (v *authenticationService) SignOut(ctx context.Context, req *pb.SignOutRequest, meta *utils.ClientMetadata) (*gp.Empty, error) {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		var authorizationTokenId uuid.UUID
		if req.AuthorizationTokenId != "" {
			authorizationTokenId = uuid.MustParse(req.AuthorizationTokenId)
		}
		jwtTokenAuthorization := &datasource.JsonWebTokenMetadata{Token: meta.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtTokenAuthorization.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		if req.All {
			var refreshTokenIds []uuid.UUID
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshTokenDeviceIdNotEqual(tx, &models.RefreshToken{DeviceId: authorizationTokenRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			for _, e := range *deleteRefreshTokenRes {
				refreshTokenIds = append(refreshTokenIds, *e.ID)
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationTokenByRefreshTokenIds(ctx, tx, &refreshTokenIds)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else if req.AuthorizationTokenId != "" {
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: &authorizationTokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{ID: authorizationTokenByReqRes.RefreshTokenId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: authorizationTokenByReqRes.ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserId: authorizationTokenRes.UserId, DeviceId: authorizationTokenRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: authorizationTokenRes.ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		}
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
}

func (v *authenticationService) ListSession(ctx context.Context, meta *utils.ClientMetadata) (*pb.ListSessionResponse, error) {
	var listSessionRes *[]models.Session
	var authorizationTokenRes *models.AuthorizationToken
	var authorizationTokenErr error
	var listSessionErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: meta.Authorization}
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
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listSessionRes, listSessionErr = v.dao.NewSessionQuery().ListSession(tx, &models.Session{UserId: authorizationTokenRes.UserId}, nil)
		if listSessionErr != nil {
			return listSessionErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var actualSession pb.Session
	otherSessions := make([]*pb.Session, 0, len(*listSessionRes))
	for _, e := range *listSessionRes {
		if *e.DeviceId != *authorizationTokenRes.DeviceId {
			otherSessions = append(otherSessions, &pb.Session{
				Id:            e.ID.String(),
				Platform:      *utils.ParsePlatformType(&e.Platform),
				SystemVersion: e.SystemVersion,
				Model:         e.Model,
				App:           *utils.ParseAppType(&e.App),
				AppVersion:    e.AppVersion,
				DeviceId:      e.DeviceId.String(),
			})
		} else {
			actualSession = pb.Session{
				Id:            e.ID.String(),
				Platform:      *utils.ParsePlatformType(&e.Platform),
				SystemVersion: e.SystemVersion,
				Model:         e.Model,
				App:           *utils.ParseAppType(&e.App),
				AppVersion:    e.AppVersion,
				DeviceId:      e.DeviceId.String(),
			}
		}
	}
	return &pb.ListSessionResponse{OtherSessions: otherSessions, ActualSession: &actualSession}, nil
}

func (v *authenticationService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, meta *utils.ClientMetadata) (*pb.RefreshTokenResponse, error) {
	var jwtAuthorizationTokenErr, jwtRefreshTokenErr error
	var jwtRefreshTokenNew, jwtAuthorizationTokenNew *datasource.JsonWebTokenMetadata
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		deviceRes, deviceErr := v.dao.NewDeviceQuery().GetDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if *deviceRes == (models.Device{}) {
			deviceRes, deviceErr = v.dao.NewDeviceQuery().CreateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier, Platform: *meta.Platform, SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId, Model: *meta.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if *deviceRes != (models.Device{}) {
			_, deviceErr := v.dao.NewDeviceQuery().UpdateDevice(tx, &models.Device{DeviceIdentifier: *meta.DeviceIdentifier}, &models.Device{SystemVersion: *meta.SystemVersion, FirebaseCloudMessagingId: *meta.FirebaseCloudMessagingId})
			if deviceErr != nil {
				return deviceErr
			}
		}
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{Token: &req.RefreshToken}
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
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: jwtRefreshToken.TokenId}, &[]string{"id"})
		if refreshTokenErr != nil {
			return refreshTokenErr
		} else if refreshTokenRes == nil {
			return errors.New("unauthenticated")
		}
		userRes, userErr := v.dao.NewUserQuery().GetUser(tx, &models.User{ID: refreshTokenRes.UserId}, &[]string{"id"})
		if userErr != nil {
			return userErr
		} else if userRes == nil {
			return errors.New("user not found")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{ID: refreshTokenRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(ctx, tx, &models.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(ctx, tx, &models.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: userRes.ID, DeviceId: deviceRes.ID, App: meta.App, AppVersion: meta.AppVersion})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		jwtRefreshTokenNew = &datasource.JsonWebTokenMetadata{TokenId: refreshTokenRes.ID}
		jwtAuthorizationTokenNew = &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
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
	return &pb.RefreshTokenResponse{RefreshToken: *jwtRefreshTokenNew.Token, AuthorizationToken: *jwtAuthorizationTokenNew.Token}, nil
}
