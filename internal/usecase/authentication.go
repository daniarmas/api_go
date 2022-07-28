package usecase

import (
	"context"
	"fmt"

	"errors"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	smtp "github.com/daniarmas/api_go/utils/smtp"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	SignIn(ctx context.Context, req *pb.SignInRequest, md *utils.ClientMetadata) (*pb.SignInResponse, error)
	SignUp(ctx context.Context, req *pb.SignUpRequest, md *utils.ClientMetadata) (*pb.SignUpResponse, error)
	SignOut(ctx context.Context, req *pb.SignOutRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	CheckSession(ctx context.Context, md *utils.ClientMetadata) (*[]string, error)
	ListSession(ctx context.Context, md *utils.ClientMetadata) (*pb.ListSessionResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, md *utils.ClientMetadata) (*pb.RefreshTokenResponse, error)
}

type authenticationService struct {
	dao    repository.Repository
	config *config.Config
	sqldb  *sqldb.Sql
}

func NewAuthenticationService(dao repository.Repository, config *config.Config, sqldb *sqldb.Sql) AuthenticationService {
	return &authenticationService{dao: dao, config: config, sqldb: sqldb}
}

func (v *authenticationService) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		user, err := v.dao.NewUserRepository().GetUser(tx, &entity.User{Email: req.Email}, &[]string{"id"})
		if err != nil {
			if err.Error() == "record not found" && (req.Type.String() == "SignIn") {
				return errors.New("user not found")
			}
		} else if user != nil && (req.Type.String() == "SignUp" || req.Type.String() == "ChangeUserEmail") {
			return errors.New("user already exists")
		}
		v.dao.NewVerificationCodeRepository().DeleteVerificationCode(tx, &entity.VerificationCode{Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *md.DeviceIdentifier}, nil)
		createVerificationCodeRes, createVerificationCodeErr := v.dao.NewVerificationCodeRepository().CreateVerificationCode(tx, &entity.VerificationCode{Code: utils.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceIdentifier: *md.DeviceIdentifier, CreateTime: time.Now(), UpdateTime: time.Now()})
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

func (v *authenticationService) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		_, err := v.dao.NewVerificationCodeRepository().GetVerificationCode(tx, &entity.VerificationCode{Code: req.Code, Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
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
	var verificationCodeRes *entity.VerificationCode
	var userRes *entity.User
	var cartItems *[]entity.CartItem
	var configuration *entity.UserConfiguration
	var deviceRes *entity.Device
	var verificationCodeErr, userErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr error
	var refreshTokenRes *entity.RefreshToken
	var authorizationTokenRes *entity.AuthorizationToken
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		deviceRes, deviceErr = v.dao.NewDeviceRepository().GetDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceRepository().CreateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else {
			_, deviceErr := v.dao.NewDeviceRepository().UpdateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		}
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeRepository().GetVerificationCode(tx, &entity.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignIn"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeRes == nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserRepository().GetUserWithAddress(tx, &entity.User{Email: req.Email}, nil)
		if userErr != nil {
			switch userErr.Error() {
			case "record not found":
				return errors.New("user not found")
			default:
				return userErr
			}
		}
		_, err := v.dao.NewVerificationCodeRepository().DeleteVerificationCode(tx, &entity.VerificationCode{Email: req.Email, Type: "SignIn", DeviceIdentifier: *md.DeviceIdentifier}, nil)
		if err != nil {
			return err
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(tx, &entity.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID}, nil)
		if deleteRefreshTokenErr != nil && deleteRefreshTokenErr.Error() != "record not found" {
			return deleteRefreshTokenErr
		}
		if deleteRefreshTokenRes != nil && len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenRepository().CreateRefreshToken(tx, &entity.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: userRes.ID, DeviceId: deviceRes.ID, App: md.App, AppVersion: md.AppVersion})
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
		cartItems, err = v.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		configuration, err = v.dao.NewUserConfigurationRepository().GetUserConfiguration(tx, &entity.UserConfiguration{UserId: userRes.ID})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.CartItem, 0, len(*cartItems))
	for _, item := range *cartItems {
		itemsResponse = append(itemsResponse, &pb.CartItem{
			Id:                   item.ID.String(),
			Name:                 item.Name,
			PriceCup:             item.PriceCup,
			ItemId:               item.ItemId.String(),
			BusinessId:           item.BusinessId.String(),
			AuthorizationTokenId: item.AuthorizationTokenId.String(),
			Quantity:             item.Quantity,
			Thumbnail:            item.Thumbnail,
			ThumbnailUrl:         v.config.ItemsBulkName + "/" + item.Thumbnail,
			BlurHash:             item.BlurHash,
			CreateTime:           timestamppb.New(item.CreateTime),
			UpdateTime:           timestamppb.New(item.UpdateTime),
		})
	}
	userAddress := make([]*pb.UserAddress, 0, len(userRes.UserAddress))
	permissions := make([]*pb.UserPermission, 0, len(userRes.UserPermissions))
	if *md.App == "Business" {
		for _, item := range userRes.UserPermissions {
			permissions = append(permissions, &pb.UserPermission{
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
			Name:           item.Name,
			Selected:       item.Selected,
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
		CartItems:           itemsResponse,
		Configuration: &pb.UserConfiguration{
			Id:                    configuration.ID.String(),
			DataSaving:            *configuration.DataSaving,
			HighQualityImagesWifi: *configuration.HighQualityImagesWifi,
			HighQualityImagesData: *configuration.HighQualityImagesData,
			UserId:                configuration.UserId.String(),
			PaymentMethod:         *utils.ParsePaymentMethodType(&configuration.PaymentMethod),
			CreateTime:            timestamppb.New(configuration.CreateTime),
			UpdateTime:            timestamppb.New(configuration.UpdateTime),
		},
		CreateTime: timestamppb.New(userRes.CreateTime),
		UpdateTime: timestamppb.New(userRes.UpdateTime),
	}}, nil
}

func (v *authenticationService) SignUp(ctx context.Context, req *pb.SignUpRequest, md *utils.ClientMetadata) (*pb.SignUpResponse, error) {
	var userRes *entity.User
	var deviceRes *entity.Device
	var verificationCodeRes *entity.VerificationCode
	var createUserAddress *entity.UserAddress
	var verificationCodeErr, userErr, deviceErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr, createUserErr error
	var refreshTokenRes *entity.RefreshToken
	var authorizationTokenRes *entity.AuthorizationToken
	var createUserRes *entity.User
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		verificationCodeRes, verificationCodeErr = v.dao.NewVerificationCodeRepository().GetVerificationCode(tx, &entity.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignUp"}, &[]string{"id"})
		if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if verificationCodeErr != nil {
			return verificationCodeErr
		}
		userRes, userErr = v.dao.NewUserRepository().GetUser(tx, &entity.User{Email: req.Email}, &[]string{"id"})
		if userErr != nil && userErr.Error() != "record not found" {
			return userErr
		} else if userRes != nil {
			return errors.New("user exists")
		}
		_, err := v.dao.NewVerificationCodeRepository().DeleteVerificationCode(tx, &entity.VerificationCode{ID: verificationCodeRes.ID}, nil)
		if err != nil {
			return err
		}
		deviceRes, deviceErr = v.dao.NewDeviceRepository().GetDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil && deviceErr.Error() != "record not found" {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceRepository().CreateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceRepository().UpdateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		}
		trueValue := true
		falseValue := false
		municipalityId := uuid.MustParse(req.UserAddress.MunicipalityId)
		provinceId := uuid.MustParse(req.UserAddress.ProvinceId)
		coordinates := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		createUserRes, createUserErr = v.dao.NewUserRepository().CreateUser(tx, &entity.User{Email: req.Email, IsLegalAge: true, FullName: req.FullName, UserConfiguration: entity.UserConfiguration{PaymentMethod: "PaymentMethodTypeCash", DataSaving: &falseValue, HighQualityImagesWifi: &trueValue, HighQualityImagesData: &trueValue}})
		if createUserErr != nil {
			return createUserErr
		}
		createUserAddress, err = v.dao.NewUserAddressRepository().CreateUserAddress(tx, &entity.UserAddress{Selected: true, Name: req.UserAddress.Name, Address: req.UserAddress.Address, Number: req.UserAddress.Number, Instructions: req.UserAddress.Instructions, ProvinceId: &provinceId, MunicipalityId: &municipalityId, Coordinates: coordinates, UserId: createUserRes.ID})
		if err != nil {
			return err
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenRepository().CreateRefreshToken(tx, &entity.RefreshToken{UserId: createUserRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: createUserRes.ID, DeviceId: deviceRes.ID, App: md.App, AppVersion: md.AppVersion})
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
			_, createBusinessUserErr := v.dao.NewBusinessUserRepository().CreateBusinessUser(tx, &entity.BusinessUser{IsBusinessOwner: isBusinessOwner, UserId: createUserRes.ID})
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
		Id:       createUserRes.ID.String(),
		FullName: createUserRes.FullName,
		Email:    createUserRes.Email,
		UserAddress: []*pb.UserAddress{
			{
				Id:             createUserAddress.ID.String(),
				Name:           createUserAddress.Name,
				Selected:       createUserAddress.Selected,
				Number:         createUserAddress.Number,
				Address:        createUserAddress.Address,
				Instructions:   createUserAddress.Instructions,
				UserId:         createUserAddress.UserId.String(),
				ProvinceId:     createUserAddress.ProvinceId.String(),
				MunicipalityId: createUserAddress.MunicipalityId.String(),
				CreateTime:     timestamppb.New(createUserAddress.CreateTime),
				UpdateTime:     timestamppb.New(createUserAddress.UpdateTime),
				Coordinates:    &pb.Point{Latitude: createUserAddress.Coordinates.FlatCoords()[1], Longitude: createUserAddress.Coordinates.FlatCoords()[0]},
			},
		},
		Configuration: &pb.UserConfiguration{
			Id:                    createUserRes.UserConfiguration.ID.String(),
			DataSaving:            *createUserRes.UserConfiguration.DataSaving,
			HighQualityImagesWifi: *createUserRes.UserConfiguration.HighQualityImagesWifi,
			HighQualityImagesData: *createUserRes.UserConfiguration.HighQualityImagesData,
			UserId:                createUserRes.UserConfiguration.UserId.String(),
			PaymentMethod:         *utils.ParsePaymentMethodType(&createUserRes.UserConfiguration.PaymentMethod),
			CreateTime:            timestamppb.New(createUserRes.UserConfiguration.CreateTime),
			UpdateTime:            timestamppb.New(createUserRes.UserConfiguration.UpdateTime),
		},
		CreateTime: timestamppb.New(createUserRes.CreateTime),
		UpdateTime: timestamppb.New(createUserRes.UpdateTime),
	}}, nil
}

func (v *authenticationService) CheckSession(ctx context.Context, md *utils.ClientMetadata) (*[]string, error) {
	var userRes *entity.User
	var deviceRes *entity.Device
	var deviceErr, userErr error
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		deviceRes, deviceErr = v.dao.NewDeviceRepository().GetDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if deviceRes == nil {
			deviceRes, deviceErr = v.dao.NewDeviceRepository().CreateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if deviceRes != nil {
			_, deviceErr := v.dao.NewDeviceRepository().UpdateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId})
			if deviceErr != nil {
				return deviceErr
			}
		}
		if md.Authorization != nil && *md.Authorization != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenRepository().GetRefreshToken(tx, &entity.RefreshToken{ID: authorizationTokenRes.RefreshTokenId}, &[]string{"id"})
			if refreshTokenErr != nil {
				return refreshTokenErr
			} else if refreshTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr = v.dao.NewUserRepository().GetUser(tx, &entity.User{ID: authorizationTokenRes.UserId}, &[]string{"id"})
			if userErr != nil {
				return userErr
			} else if userRes == nil {
				return errors.New("user not found")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &[]string{}, nil
}

func (v *authenticationService) SignOut(ctx context.Context, req *pb.SignOutRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		var authorizationTokenId uuid.UUID
		if req.AuthorizationTokenId != "" {
			authorizationTokenId = uuid.MustParse(req.AuthorizationTokenId)
		}
		jwtTokenAuthorization := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtTokenAuthorization.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		if req.All {
			var refreshTokenIds []uuid.UUID
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshTokenDeviceIdNotEqual(tx, &entity.RefreshToken{DeviceId: authorizationTokenRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			for _, e := range *deleteRefreshTokenRes {
				refreshTokenIds = append(refreshTokenIds, *e.ID)
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationTokenByRefreshTokenIds(ctx, tx, &refreshTokenIds)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else if req.AuthorizationTokenId != "" {
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: &authorizationTokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(tx, &entity.RefreshToken{ID: authorizationTokenByReqRes.RefreshTokenId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: authorizationTokenByReqRes.ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(tx, &entity.RefreshToken{UserId: authorizationTokenRes.UserId, DeviceId: authorizationTokenRes.DeviceId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: authorizationTokenRes.ID}, nil)
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

func (v *authenticationService) ListSession(ctx context.Context, md *utils.ClientMetadata) (*pb.ListSessionResponse, error) {
	var listSessionRes *[]entity.Session
	var authorizationTokenRes *entity.AuthorizationToken
	var authorizationTokenErr error
	var listSessionErr error
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listSessionRes, listSessionErr = v.dao.NewSessionRepository().ListSession(tx, &entity.Session{UserId: authorizationTokenRes.UserId}, nil)
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

func (v *authenticationService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, md *utils.ClientMetadata) (*pb.RefreshTokenResponse, error) {
	var jwtAuthorizationTokenErr, jwtRefreshTokenErr error
	var jwtRefreshTokenNew, jwtAuthorizationTokenNew *datasource.JsonWebTokenMetadata
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		appErr := v.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		deviceRes, deviceErr := v.dao.NewDeviceRepository().GetDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &[]string{"id"})
		if deviceErr != nil {
			return deviceErr
		} else if *deviceRes == (entity.Device{}) {
			deviceRes, deviceErr = v.dao.NewDeviceRepository().CreateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if deviceErr != nil {
				return deviceErr
			}
		} else if *deviceRes != (entity.Device{}) {
			_, deviceErr := v.dao.NewDeviceRepository().UpdateDevice(tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId})
			if deviceErr != nil {
				return deviceErr
			}
		}
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		refreshTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtRefreshToken(jwtRefreshToken)
		if refreshTokenParseErr != nil {
			switch refreshTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("refresh token expired")
			case "signature is invalid":
				return errors.New("refresh token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("refresh token contains an invalid number of segments")
			default:
				return refreshTokenParseErr
			}
		}
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenRepository().GetRefreshToken(tx, &entity.RefreshToken{ID: jwtRefreshToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if refreshTokenErr != nil && refreshTokenErr.Error() == "record not found" {
			return errors.New("refresh token not found")
		} else if refreshTokenErr != nil {
			return refreshTokenErr
		}
		userRes, userErr := v.dao.NewUserRepository().GetUser(tx, &entity.User{ID: refreshTokenRes.UserId}, &[]string{"id"})
		if userErr != nil {
			return userErr
		} else if userRes == nil {
			return errors.New("user not found")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(tx, &entity.RefreshToken{ID: refreshTokenRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenRepository().CreateRefreshToken(tx, &entity.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: userRes.ID, DeviceId: deviceRes.ID, App: md.App, AppVersion: md.AppVersion})
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