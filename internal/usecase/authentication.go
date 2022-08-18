package usecase

import (
	"context"
	"fmt"

	"errors"
	"time"

	"firebase.google.com/go/messaging"
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
	CheckSession(ctx context.Context, md *utils.ClientMetadata) (*pb.CheckSessionResponse, error)
	ListSession(ctx context.Context, md *utils.ClientMetadata) (*pb.ListSessionResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, md *utils.ClientMetadata) (*pb.RefreshTokenResponse, error)
	SessionExists(ctx context.Context, req *pb.SessionExistsRequest, md *utils.ClientMetadata) (*pb.SessionExistsResponse, error)
}

type authenticationService struct {
	dao                repository.Repository
	config             *config.Config
	sqldb              *sqldb.Sql
	moolShoppingClient *messaging.Client
}

func NewAuthenticationService(dao repository.Repository, config *config.Config, sqldb *sqldb.Sql, moolShoppingClient *messaging.Client) AuthenticationService {
	return &authenticationService{dao: dao, config: config, sqldb: sqldb, moolShoppingClient: moolShoppingClient}
}

func (v *authenticationService) SessionExists(ctx context.Context, req *pb.SessionExistsRequest, md *utils.ClientMetadata) (*pb.SessionExistsResponse, error) {
	var res pb.SessionExistsResponse
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		_, err = v.dao.NewVerificationCodeRepository().GetVerificationCode(ctx, tx, &entity.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignIn"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("verification code not found")
		}
		user, err := v.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{Email: req.Email})
		if err != nil && err.Error() == "record not found" {
			return errors.New("user not found")
		} else if err != nil {
			return err
		}
		authorizationToken, err := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{UserId: user.ID})
		if err != nil && err.Error() == "record not found" {
			return errors.New("session not exists")
		} else if err != nil {
			return err
		}
		res.AuthorizationTokenId = authorizationToken.ID.String()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (v *authenticationService) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		device, err := v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if device == nil {
			_, err = v.dao.NewDeviceRepository().CreateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		} else {
			_, err = v.dao.NewDeviceRepository().UpdateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		}
		_, err = v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		user, err := v.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{Email: req.Email})
		if err != nil {
			if err.Error() == "record not found" && (req.Type.String() == "SignIn") {
				return errors.New("user not found")
			}
		} else if user != nil && (req.Type.String() == "SignUp" || req.Type.String() == "ChangeUserEmail") {
			return errors.New("user already exists")
		}
		v.dao.NewVerificationCodeRepository().DeleteVerificationCode(ctx, tx, &entity.VerificationCode{Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *md.DeviceIdentifier}, nil)
		createVerificationCodeRes, err := v.dao.NewVerificationCodeRepository().CreateVerificationCode(ctx, tx, &entity.VerificationCode{Code: utils.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceIdentifier: *md.DeviceIdentifier, CreateTime: time.Now(), UpdateTime: time.Now()})
		if err != nil {
			return err
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
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		_, err = v.dao.NewVerificationCodeRepository().GetVerificationCode(ctx, tx, &entity.VerificationCode{Code: req.Code, Email: req.Email, Type: req.Type.String(), DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if err != nil {
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
	var user *entity.User
	var cartItems *[]entity.CartItem
	var configuration *entity.UserConfiguration
	var app *entity.Application
	var authToken *entity.AuthorizationToken
	// var actualDevice *entity.Device
	// var signOutDevice *entity.Device
	var existsUpcomingOrders *bool
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		actualDevice, err := v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if actualDevice == nil {
			actualDevice, err = v.dao.NewDeviceRepository().CreateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		} else {
			_, err = v.dao.NewDeviceRepository().UpdateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		}
		app, err = v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		_, err = v.dao.NewVerificationCodeRepository().GetVerificationCode(ctx, tx, &entity.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignIn"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if err != nil {
			return err
		}
		user, err = v.dao.NewUserRepository().GetUserWithAddress(ctx, tx, &entity.User{Email: req.Email})
		if err != nil {
			switch err.Error() {
			case "record not found":
				return errors.New("user not found")
			default:
				return err
			}
		}
		// Limit session to one by device
		authToken, err = v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{UserId: user.ID})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if authToken != nil && !req.Logout {
			return errors.New("session limit reached")
		} else if authToken != nil && req.Logout {
			deleteRefreshTokenRes, err := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(ctx, tx, &entity.RefreshToken{UserId: user.ID}, nil)
			if err != nil && err.Error() != "record not found" {
				return err
			}
			if deleteRefreshTokenRes != nil && len(*deleteRefreshTokenRes) != 0 {
				_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
				if deleteAuthorizationTokenErr != nil {
					return deleteAuthorizationTokenErr
				}
				// if authToken != nil && req.Logout {
				// deviceId := *deleteAuthorizationTokenRes
				// signOutDevice, err = v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{ID: deviceId[0].ID})
				// if err != nil && err.Error() != "record not found" {
				// 	return err
				// }
				// }
			}
		}
		//
		_, err = v.dao.NewVerificationCodeRepository().DeleteVerificationCode(ctx, tx, &entity.VerificationCode{Email: req.Email, Type: "SignIn", DeviceIdentifier: *md.DeviceIdentifier}, nil)
		if err != nil {
			return err
		}
		deleteRefreshTokenRes, err := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(ctx, tx, &entity.RefreshToken{UserId: user.ID, DeviceId: actualDevice.ID}, nil)
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if deleteRefreshTokenRes != nil && len(*deleteRefreshTokenRes) != 0 {
			_, err = v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
			if err != nil {
				return err
			}
		}
		createRefreshTokenRes, err := v.dao.NewRefreshTokenRepository().CreateRefreshToken(ctx, tx, &entity.RefreshToken{UserId: user.ID, DeviceId: actualDevice.ID})
		if err != nil {
			return err
		}
		createAuthTokenRes, err := v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: createRefreshTokenRes.ID, UserId: user.ID, DeviceId: actualDevice.ID, App: &app.Name, AppVersion: &app.Version})
		if err != nil {
			return err
		}
		jwtRefreshToken = &datasource.JsonWebTokenMetadata{TokenId: createRefreshTokenRes.ID}
		err = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken = &datasource.JsonWebTokenMetadata{TokenId: createAuthTokenRes.ID}
		err = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			return err
		}
		cartItems, err = v.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: createAuthTokenRes.UserId})
		if err != nil {
			return err
		}
		configuration, err = v.dao.NewUserConfigurationRepository().GetUserConfiguration(ctx, tx, &entity.UserConfiguration{UserId: user.ID})
		if err != nil {
			return err
		}
		existsUpcomingOrders, err = v.dao.NewOrderRepository().ExistsUpcomingOrders(tx, *user.ID)
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
	userAddress := make([]*pb.UserAddress, 0, len(user.UserAddress))
	permissions := make([]*pb.UserPermission, 0, len(user.UserPermissions))
	if app.Name == "Mool for business" {
		for _, item := range user.UserPermissions {
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
	for _, item := range user.UserAddress {
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
	var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
	if user.HighQualityPhoto != "" {
		highQualityPhotoUrl = v.config.UsersBulkName + "/" + user.HighQualityPhoto
		lowQualityPhotoUrl = v.config.UsersBulkName + "/" + user.LowQualityPhoto
		thumbnailUrl = v.config.UsersBulkName + "/" + user.Thumbnail

	}
	// if authToken != nil && req.Logout {
	// 	bodyMsg := fmt.Sprintf("Se ha iniciado una nueva sesión en un dispositivo " + actualDevice.Model)
	// 	// See documentation on defining a message payload.
	// 	message := &messaging.Message{
	// 		Notification: &messaging.Notification{
	// 			Title: "Inició de sesión",
	// 			Body:  bodyMsg,
	// 		},
	// 		Data: map[string]string{
	// 			"res_id": "",
	// 			"res":    "session",
	// 			"verb":   "create",
	// 		},
	// 		Token: signOutDevice.FirebaseCloudMessagingId,
	// 	}

	// 	// Send a message to the device corresponding to the provided
	// 	// registration token.
	// 	response, err := v.moolShoppingClient.Send(ctx, message)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// Response is a message ID string.
	// 	fmt.Println("Successfully sent message:", response)
	// }
	go smtp.SendSignInMail(req.Email, time.Now(), v.config, md)
	return &pb.SignInResponse{AuthorizationToken: *jwtAuthorizationToken.Token, RefreshToken: *jwtRefreshToken.Token, User: &pb.User{
		Id:                   user.ID.String(),
		FullName:             user.FullName,
		Email:                user.Email,
		HighQualityPhoto:     user.HighQualityPhoto,
		HighQualityPhotoUrl:  highQualityPhotoUrl,
		LowQualityPhoto:      user.LowQualityPhoto,
		LowQualityPhotoUrl:   lowQualityPhotoUrl,
		Thumbnail:            user.Thumbnail,
		ThumbnailUrl:         thumbnailUrl,
		BlurHash:             user.BlurHash,
		Permissions:          permissions,
		UserAddress:          userAddress,
		CartItems:            itemsResponse,
		ExistsUpcomingOrders: *existsUpcomingOrders,
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
		CreateTime: timestamppb.New(user.CreateTime),
		UpdateTime: timestamppb.New(user.UpdateTime),
	}}, nil
}

func (v *authenticationService) SignUp(ctx context.Context, req *pb.SignUpRequest, md *utils.ClientMetadata) (*pb.SignUpResponse, error) {
	var getUserRes *entity.User
	var getDeviceRes *entity.Device
	var getVerificationCode *entity.VerificationCode
	var createUserAddressRes *entity.UserAddress
	var createUserRes *entity.User
	var (
		jwtRefreshToken       *datasource.JsonWebTokenMetadata
		jwtAuthorizationToken *datasource.JsonWebTokenMetadata
	)
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		app, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		getVerificationCode, err = v.dao.NewVerificationCodeRepository().GetVerificationCode(ctx, tx, &entity.VerificationCode{Email: req.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "SignUp"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("verification code not found")
		} else if err != nil {
			return err
		}
		getUserRes, err = v.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{Email: req.Email})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if getUserRes != nil {
			return errors.New("user exists")
		}
		_, err = v.dao.NewVerificationCodeRepository().DeleteVerificationCode(ctx, tx, &entity.VerificationCode{ID: getVerificationCode.ID}, nil)
		if err != nil {
			return err
		}
		getDeviceRes, err = v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if getDeviceRes == nil {
			getDeviceRes, err = v.dao.NewDeviceRepository().CreateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		} else if getDeviceRes != nil {
			_, err := v.dao.NewDeviceRepository().UpdateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		}
		trueValue := true
		falseValue := false
		municipalityId := uuid.MustParse(req.UserAddress.MunicipalityId)
		provinceId := uuid.MustParse(req.UserAddress.ProvinceId)
		coordinates := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		createUserRes, err = v.dao.NewUserRepository().CreateUser(ctx, tx, &entity.User{Email: req.Email, IsLegalAge: true, FullName: req.FullName, UserConfiguration: entity.UserConfiguration{PaymentMethod: "PaymentMethodTypeCash", DataSaving: &falseValue, HighQualityImagesWifi: &trueValue, HighQualityImagesData: &trueValue}})
		if err != nil {
			return err
		}
		createUserAddressRes, err = v.dao.NewUserAddressRepository().CreateUserAddress(tx, &entity.UserAddress{Selected: true, Name: req.UserAddress.Name, Address: req.UserAddress.Address, Number: req.UserAddress.Number, Instructions: req.UserAddress.Instructions, ProvinceId: &provinceId, MunicipalityId: &municipalityId, Coordinates: coordinates, UserId: createUserRes.ID})
		if err != nil {
			return err
		}
		refreshTokenRes, err := v.dao.NewRefreshTokenRepository().CreateRefreshToken(ctx, tx, &entity.RefreshToken{UserId: createUserRes.ID, DeviceId: getDeviceRes.ID})
		if err != nil {
			return err
		}
		authorizationTokenRes, err := v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: createUserRes.ID, DeviceId: getDeviceRes.ID, App: &app.Name, AppVersion: &app.Version})
		if err != nil {
			return err
		}
		jwtRefreshToken = &datasource.JsonWebTokenMetadata{TokenId: refreshTokenRes.ID}
		err = repository.Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(jwtRefreshToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken = &datasource.JsonWebTokenMetadata{TokenId: authorizationTokenRes.ID}
		err = repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			return err
		}
		var isBusinessOwner bool = false
		if getVerificationCode.Type == "SignUpBusinessOwner" {
			isBusinessOwner = true
		}
		if req.SignUpType.String() == "SignUpBusiness" {
			_, createBusinessUserErr := v.dao.NewBusinessUserRepository().CreateBusinessUser(ctx, tx, &entity.BusinessUser{IsBusinessOwner: isBusinessOwner, UserId: createUserRes.ID})
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
				Id:             createUserAddressRes.ID.String(),
				Name:           createUserAddressRes.Name,
				Selected:       createUserAddressRes.Selected,
				Number:         createUserAddressRes.Number,
				Address:        createUserAddressRes.Address,
				Instructions:   createUserAddressRes.Instructions,
				UserId:         createUserAddressRes.UserId.String(),
				ProvinceId:     createUserAddressRes.ProvinceId.String(),
				MunicipalityId: createUserAddressRes.MunicipalityId.String(),
				CreateTime:     timestamppb.New(createUserAddressRes.CreateTime),
				UpdateTime:     timestamppb.New(createUserAddressRes.UpdateTime),
				Coordinates:    &pb.Point{Latitude: createUserAddressRes.Coordinates.FlatCoords()[1], Longitude: createUserAddressRes.Coordinates.FlatCoords()[0]},
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

func (v *authenticationService) CheckSession(ctx context.Context, md *utils.ClientMetadata) (*pb.CheckSessionResponse, error) {
	var res pb.CheckSessionResponse
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		deviceRes, err := v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if deviceRes == nil {
			_, err = v.dao.NewDeviceRepository().CreateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		} else if deviceRes != nil {
			_, err = v.dao.NewDeviceRepository().UpdateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId})
			if err != nil {
				return err
			}
		}
		app, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		if md.Authorization != nil && *md.Authorization != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
			authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
			if authorizationTokenParseErr != nil {
				switch authorizationTokenParseErr.Error() {
				case "Token is expired":
					return errors.New("authorization token expired")
				case "signature is invalid":
					return errors.New("authorization token signature is invalid")
				case "token contains an invalid number of segments":
					return errors.New("authorization token contains an invalid number of segments")
				default:
					return authorizationTokenParseErr
				}
			}
			authorizationTokenRes, err := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
			if err != nil && err.Error() == "record not found" {
				return errors.New("unauthenticated user")
			} else if err != nil {
				return err
			}
			userRes, err := v.dao.NewUserRepository().GetUserWithAddress(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
			if err != nil {
				return err
			} else if userRes == nil {
				return errors.New("unauthenticated user")
			}
			cartItems, err := v.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
			if err != nil {
				return err
			}
			configuration, err := v.dao.NewUserConfigurationRepository().GetUserConfiguration(ctx, tx, &entity.UserConfiguration{UserId: userRes.ID})
			if err != nil {
				return err
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
			if app.Name == "Mool for business" {
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
			var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
			if userRes.HighQualityPhoto != "" {
				highQualityPhotoUrl = v.config.UsersBulkName + "/" + userRes.HighQualityPhoto
				lowQualityPhotoUrl = v.config.UsersBulkName + "/" + userRes.LowQualityPhoto
				thumbnailUrl = v.config.UsersBulkName + "/" + userRes.Thumbnail
			}
			existsUpcomingOrders, err := v.dao.NewOrderRepository().ExistsUpcomingOrders(tx, *userRes.ID)
			if err != nil {
				return err
			}
			res.User = &pb.User{
				Id:                   userRes.ID.String(),
				FullName:             userRes.FullName,
				Email:                userRes.Email,
				HighQualityPhotoUrl:  highQualityPhotoUrl,
				HighQualityPhoto:     userRes.HighQualityPhoto,
				LowQualityPhotoUrl:   lowQualityPhotoUrl,
				LowQualityPhoto:      userRes.LowQualityPhoto,
				ThumbnailUrl:         thumbnailUrl,
				Thumbnail:            userRes.Thumbnail,
				BlurHash:             userRes.BlurHash,
				Permissions:          permissions,
				UserAddress:          userAddress,
				CartItems:            itemsResponse,
				ExistsUpcomingOrders: *existsUpcomingOrders,
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
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (v *authenticationService) SignOut(ctx context.Context, req *pb.SignOutRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		var authorizationTokenId uuid.UUID
		if req.AuthorizationTokenId != "" {
			authorizationTokenId = uuid.MustParse(req.AuthorizationTokenId)
		}
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		if req.All {
			var refreshTokenIds []uuid.UUID
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshTokenDeviceIdNotEqual(ctx, tx, &entity.RefreshToken{DeviceId: authorizationTokenRes.DeviceId}, nil)
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
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: &authorizationTokenId})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(ctx, tx, &entity.RefreshToken{ID: authorizationTokenByReqRes.RefreshTokenId}, nil)
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: authorizationTokenByReqRes.ID}, nil)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			_, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(ctx, tx, &entity.RefreshToken{UserId: authorizationTokenRes.UserId, DeviceId: authorizationTokenRes.DeviceId}, nil)
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
	var listSessionErr error
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, err := v.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		listSessionRes, listSessionErr = v.dao.NewSessionRepository().ListSession(tx, &entity.Session{UserId: authorizationTokenRes.UserId})
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
		app, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		deviceRes, err := v.dao.NewDeviceRepository().GetDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if deviceRes == nil {
			deviceRes, err = v.dao.NewDeviceRepository().CreateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		} else {
			_, err = v.dao.NewDeviceRepository().UpdateDevice(ctx, tx, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier}, &entity.Device{DeviceIdentifier: *md.DeviceIdentifier, Platform: *md.Platform, SystemVersion: *md.SystemVersion, FirebaseCloudMessagingId: *md.FirebaseCloudMessagingId, Model: *md.Model})
			if err != nil {
				return err
			}
		}
		jwtRefreshToken := &datasource.JsonWebTokenMetadata{Token: &req.RefreshToken}
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
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenRepository().GetRefreshToken(ctx, tx, &entity.RefreshToken{ID: jwtRefreshToken.TokenId})
		if refreshTokenErr != nil && refreshTokenErr.Error() == "record not found" {
			return errors.New("refresh token not found")
		} else if refreshTokenErr != nil {
			return refreshTokenErr
		}
		userRes, userErr := v.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{ID: refreshTokenRes.UserId})
		if userErr != nil {
			return userErr
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenRepository().DeleteRefreshToken(ctx, tx, &entity.RefreshToken{ID: refreshTokenRes.ID}, nil)
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenRepository().DeleteAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: (*deleteRefreshTokenRes)[0].ID}, nil)
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenRepository().CreateRefreshToken(ctx, tx, &entity.RefreshToken{UserId: userRes.ID, DeviceId: deviceRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenRepository().CreateAuthorizationToken(ctx, tx, &entity.AuthorizationToken{RefreshTokenId: refreshTokenRes.ID, UserId: userRes.ID, DeviceId: deviceRes.ID, App: &app.Name, AppVersion: &app.Version})
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
