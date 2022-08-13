package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, md *utils.ClientMetadata) (*pb.User, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, md *utils.ClientMetadata) (*pb.User, error)
	GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest, md *utils.ClientMetadata) (*pb.GetAddressInfoResponse, error)
	ListUserAddress(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.ListUserAddressResponse, error)
	GetUserAddress(ctx context.Context, req *pb.GetUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error)
	CreateUserAddress(ctx context.Context, req *pb.CreateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error)
	UpdateUserAddress(ctx context.Context, req *pb.UpdateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error)
	DeleteUserAddress(ctx context.Context, req *pb.DeleteUserAddressRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	UpdateUserConfiguration(ctx context.Context, req *pb.UpdateUserConfigurationRequest, md *utils.ClientMetadata) (*pb.UserConfiguration, error)
}

type userService struct {
	dao    repository.Repository
	config *config.Config
	rdb    *redis.Client
	sqldb  *sqldb.Sql
}

func NewUserService(dao repository.Repository, config *config.Config, rdb *redis.Client, sqldb *sqldb.Sql) UserService {
	return &userService{dao: dao, config: config, rdb: rdb, sqldb: sqldb}
}

func (i *userService) UpdateUserConfiguration(ctx context.Context, req *pb.UpdateUserConfigurationRequest, md *utils.ClientMetadata) (*pb.UserConfiguration, error) {
	var res pb.UserConfiguration
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		dataSaving := req.UserConfiguration.DataSaving
		highQualityImagesWifi := req.UserConfiguration.HighQualityImagesWifi
		highQualityImagesData := req.UserConfiguration.HighQualityImagesData
		data := entity.UserConfiguration{
			DataSaving:            &dataSaving,
			HighQualityImagesWifi: &highQualityImagesWifi,
			HighQualityImagesData: &highQualityImagesData,
			PaymentMethod:         req.UserConfiguration.PaymentMethod.String(),
		}
		if req.UserConfiguration.PaymentMethod == pb.PaymentMethodType_PaymentMethodTypeUnspecified {
			data.PaymentMethod = ""
		}
		userConfigurationRes, err := i.dao.NewUserConfigurationRepository().UpdateUserConfiguration(ctx, tx, &entity.UserConfiguration{UserId: authorizationTokenRes.UserId}, &data)
		if err != nil && err.Error() == "record not found" {
			return errors.New("user configuration not found")
		} else if err != nil {
			return err
		}
		res = pb.UserConfiguration{
			Id:                    userConfigurationRes.ID.String(),
			DataSaving:            *userConfigurationRes.DataSaving,
			HighQualityImagesWifi: *userConfigurationRes.HighQualityImagesWifi,
			HighQualityImagesData: *userConfigurationRes.HighQualityImagesData,
			PaymentMethod:         *utils.ParsePaymentMethodType(&userConfigurationRes.PaymentMethod),
			UserId:                userConfigurationRes.UserId.String(),
			CreateTime:            timestamppb.New(userConfigurationRes.CreateTime),
			UpdateTime:            timestamppb.New(userConfigurationRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) GetUserAddress(ctx context.Context, req *pb.GetUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error) {
	var res pb.UserAddress
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		userRes, err := i.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
		if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		userAddressRes, err := i.dao.NewUserAddressRepository().GetUserAddress(tx, &entity.UserAddress{UserId: userRes.ID, ID: &id})
		if err != nil {
			return err
		}
		res = pb.UserAddress{
			Id:             userAddressRes.ID.String(),
			Name:           userAddressRes.Name,
			Number:         userAddressRes.Number,
			Instructions:   userAddressRes.Instructions,
			Address:        userAddressRes.Address,
			Selected:       userAddressRes.Selected,
			UserId:         userAddressRes.UserId.String(),
			ProvinceId:     userAddressRes.ProvinceId.String(),
			MunicipalityId: userAddressRes.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: userAddressRes.Coordinates.FlatCoords()[1], Longitude: userAddressRes.Coordinates.FlatCoords()[0]},
			CreateTime:     timestamppb.New(userAddressRes.CreateTime),
			UpdateTime:     timestamppb.New(userAddressRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) UpdateUserAddress(ctx context.Context, req *pb.UpdateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error) {
	var res pb.UserAddress
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		userRes, err := i.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
		if err != nil {
			return err
		}
		listAddressRes, err := i.dao.NewUserAddressRepository().ListUserAddress(tx, &entity.UserAddress{UserId: userRes.ID})
		if err != nil {
			return err
		}
		if len(*listAddressRes) == 10 {
			return errors.New("only can have 10 user_address")
		}
		id := uuid.MustParse(req.Id)
		if req.UserAddress.Selected {
			_, err := i.dao.NewUserAddressRepository().UpdateUserAddressByUserId(tx, &entity.UserAddress{UserId: authorizationTokenRes.UserId}, &entity.UserAddress{Selected: false})
			if err != nil && err.Error() == "record not found" {
				return errors.New("user address not found")
			} else if err != nil {
				return err
			}
		}
		where := entity.UserAddress{
			Name:         req.UserAddress.Name,
			Address:      req.UserAddress.Address,
			UserId:       userRes.ID,
			Instructions: req.UserAddress.Instructions,
			Number:       req.UserAddress.Number,
			Selected:     req.UserAddress.Selected,
		}
		if req.UserAddress.ProvinceId != "" {
			value := uuid.MustParse(req.UserAddress.ProvinceId)
			where.ProvinceId = &value
		}
		if req.UserAddress.MunicipalityId != "" {
			value := uuid.MustParse(req.UserAddress.MunicipalityId)
			where.MunicipalityId = &value
		}
		if req.UserAddress.Coordinates != nil {
			where.Coordinates = ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		}
		var updateUserAddressRes *entity.UserAddress
		updateUserAddressRes, err = i.dao.NewUserAddressRepository().UpdateUserAddress(tx, &entity.UserAddress{ID: &id}, &where)
		if err != nil && err.Error() == "record not found" {
			return errors.New("user address not found")
		} else if err != nil {
			return err
		}
		res = pb.UserAddress{
			Id:             updateUserAddressRes.ID.String(),
			Name:           updateUserAddressRes.Name,
			Number:         updateUserAddressRes.Number,
			Instructions:   updateUserAddressRes.Instructions,
			Address:        updateUserAddressRes.Address,
			Selected:       updateUserAddressRes.Selected,
			UserId:         updateUserAddressRes.UserId.String(),
			ProvinceId:     updateUserAddressRes.ProvinceId.String(),
			MunicipalityId: updateUserAddressRes.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: updateUserAddressRes.Coordinates.FlatCoords()[1], Longitude: updateUserAddressRes.Coordinates.FlatCoords()[0]},
			CreateTime:     timestamppb.New(updateUserAddressRes.CreateTime),
			UpdateTime:     timestamppb.New(updateUserAddressRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) DeleteUserAddress(ctx context.Context, req *pb.DeleteUserAddressRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		_, err = i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		_, err = i.dao.NewUserAddressRepository().DeleteUserAddress(tx, &entity.UserAddress{ID: &id}, nil)
		if err != nil && err.Error() == "record not found" {
			return errors.New("user address not found")
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

func (i *userService) CreateUserAddress(ctx context.Context, req *pb.CreateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error) {
	var res pb.UserAddress
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr := i.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
		if userErr != nil {
			return userErr
		}
		listAddressRes, listAddressErr := i.dao.NewUserAddressRepository().ListUserAddress(tx, &entity.UserAddress{UserId: userRes.ID})
		if listAddressErr != nil {
			return listAddressErr
		}
		if len(*listAddressRes) == 10 {
			return errors.New("only can have 10 user_address")
		}
		provinceId := uuid.MustParse(req.UserAddress.ProvinceId)
		municipalityId := uuid.MustParse(req.UserAddress.MunicipalityId)
		location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		createUserAddressRes, createUserAddressErr := i.dao.NewUserAddressRepository().CreateUserAddress(tx, &entity.UserAddress{Selected: req.UserAddress.Selected, Name: req.UserAddress.Name, Address: req.UserAddress.Address, Number: req.UserAddress.Number, Instructions: req.UserAddress.Instructions, UserId: userRes.ID, ProvinceId: &provinceId, MunicipalityId: &municipalityId, Coordinates: location})
		if createUserAddressErr != nil {
			return createUserAddressErr
		}
		res = pb.UserAddress{
			Id:             createUserAddressRes.ID.String(),
			Name:           createUserAddressRes.Name,
			Number:         createUserAddressRes.Number,
			Instructions:   createUserAddressRes.Instructions,
			Address:        createUserAddressRes.Address,
			UserId:         createUserAddressRes.UserId.String(),
			Selected:       createUserAddressRes.Selected,
			ProvinceId:     createUserAddressRes.ProvinceId.String(),
			MunicipalityId: createUserAddressRes.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: createUserAddressRes.Coordinates.FlatCoords()[1], Longitude: createUserAddressRes.Coordinates.FlatCoords()[0]},
			CreateTime:     timestamppb.New(createUserAddressRes.CreateTime),
			UpdateTime:     timestamppb.New(createUserAddressRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) ListUserAddress(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.ListUserAddressResponse, error) {
	var res pb.ListUserAddressResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userId := *authorizationTokenRes.UserId
		listUserAddressRes, listUserAddressErr := i.dao.NewUserAddressRepository().ListUserAddress(tx, &entity.UserAddress{UserId: &userId})
		if listUserAddressErr != nil {
			return listUserAddressErr
		}
		userAddress := make([]*pb.UserAddress, 0, len(*listUserAddressRes))
		for _, i := range *listUserAddressRes {
			userAddress = append(userAddress, &pb.UserAddress{
				Id:             i.ID.String(),
				Name:           i.Name,
				UserId:         i.UserId.String(),
				Coordinates:    &pb.Point{Latitude: i.Coordinates.FlatCoords()[1], Longitude: i.Coordinates.FlatCoords()[0]},
				Address:        i.Address,
				Number:         i.Number,
				Selected:       i.Selected,
				Instructions:   i.Instructions,
				ProvinceId:     i.ProvinceId.String(),
				MunicipalityId: i.MunicipalityId.String(),
				CreateTime:     timestamppb.New(i.CreateTime),
				UpdateTime:     timestamppb.New(i.UpdateTime),
			})
		}
		res.UserAddress = userAddress
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest, md *utils.ClientMetadata) (*pb.GetAddressInfoResponse, error) {
	var res pb.GetAddressInfoResponse
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		muncipalityRes, err := i.dao.NewMunicipalityRepository().GetMunicipalityByCoordinate(tx, location)
		if err != nil && err.Error() == "record not found" {
			return errors.New("municipality not found")
		} else if err != nil {
			return err
		}
		provinceRes, err := i.dao.NewProvinceRepository().GetProvince(tx, &entity.Province{ID: muncipalityRes.ProvinceId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("province not found")
		} else if err != nil {
			return err
		}
		res = pb.GetAddressInfoResponse{ProvinceId: provinceRes.ID.String(), ProvinceName: provinceRes.Name, ProvinceNameAbbreviation: provinceRes.Codename, MunicipalityName: muncipalityRes.Name, MunicipalityId: muncipalityRes.ID.String()}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *userService) GetUser(ctx context.Context, md *utils.ClientMetadata) (*pb.User, error) {
	var userRes *entity.User
	var userErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr = i.dao.NewUserRepository().GetUserWithAddress(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	userAddress := make([]*pb.UserAddress, 0, len(userRes.UserAddress))
	permissions := make([]*pb.UserPermission, 0, len(userRes.UserPermissions))
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
	for _, item := range userRes.UserAddress {
		userAddress = append(userAddress, &pb.UserAddress{
			Id:             item.ID.String(),
			Name:           item.Name,
			Number:         item.Number,
			Address:        item.Address,
			Instructions:   item.Instructions,
			Selected:       item.Selected,
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
		highQualityPhotoUrl = i.config.UsersBulkName + "/" + userRes.HighQualityPhoto
		lowQualityPhotoUrl = i.config.UsersBulkName + "/" + userRes.LowQualityPhoto
		thumbnailUrl = i.config.UsersBulkName + "/" + userRes.Thumbnail

	}
	return &pb.User{
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
		UserAddress:         userAddress,
		Permissions:         permissions,
		CreateTime:          timestamppb.New(userRes.CreateTime),
		UpdateTime:          timestamppb.New(userRes.UpdateTime),
	}, nil
}

func (i *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, md *utils.ClientMetadata) (*pb.User, error) {
	var updatedUserRes *entity.User
	var updatedUserErr error
	var userId uuid.UUID
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		// chech if is the user or if have permission
		if req.User.Id != "" && authorizationTokenRes.UserId.String() != req.User.Id {
			_, err := i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{Name: "admin"})
			if err != nil && err.Error() == "record not found" {
				return errors.New("not have permission")
			} else if err != nil && err.Error() != "record not found" {
				return err
			}
		}
		userRes, userErr := i.dao.NewUserRepository().GetUser(ctx, tx, &entity.User{ID: authorizationTokenRes.UserId})
		if userErr != nil {
			return userErr
		}
		if req.User.Id == "" {
			userId = *userRes.ID
		} else {
			userId = uuid.MustParse(req.User.Id)
		}
		if req.User.Email != "" {
			if req.Code == "" {
				return errors.New("missing code")
			}
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeRepository().GetVerificationCode(ctx, tx, &entity.VerificationCode{Email: userRes.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "ChangeUserEmail"})
			if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
				return errors.New("verification code not found")
			} else if verificationCodeErr != nil {
				return verificationCodeErr
			}
			_, err := i.dao.NewVerificationCodeRepository().DeleteVerificationCode(ctx, tx, &entity.VerificationCode{ID: verificationCode.ID}, nil)
			if err != nil {
				return err
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(ctx, tx, &entity.User{ID: &userId}, &entity.User{Email: req.User.Email})
			if updatedUserErr != nil {
				return updatedUserErr
			}

		} else if req.User.HighQualityPhoto != "" && req.User.LowQualityPhoto != "" && req.User.Thumbnail != "" && req.User.BlurHash != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.UsersBulkName, req.User.HighQualityPhoto)
			if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.UsersBulkName, req.User.LowQualityPhoto)
			if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.UsersBulkName, req.User.Thumbnail)
			if tnErr != nil {
				return tnErr
			}
			if userRes.HighQualityPhoto != "" {
				_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.UsersDeletedBulkName, Object: userRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.UsersBulkName, Object: userRes.HighQualityPhoto})
				if copyHqErr != nil {
					return copyHqErr
				}
			}
			if userRes.LowQualityPhoto != "" {
				_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.UsersDeletedBulkName, Object: userRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.UsersBulkName, Object: userRes.LowQualityPhoto})
				if copyLqErr != nil {
					return copyLqErr
				}
			}
			if userRes.Thumbnail != "" {
				_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.UsersDeletedBulkName, Object: userRes.Thumbnail}, minio.CopySrcOptions{Bucket: i.config.UsersBulkName, Object: userRes.Thumbnail})
				if copyThErr != nil {
					return copyThErr
				}
			}
			if userRes.HighQualityPhoto != "" {
				rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.UsersBulkName, userRes.HighQualityPhoto, minio.RemoveObjectOptions{})
				if rmHqErr != nil {
					return rmHqErr
				}
			}
			if userRes.LowQualityPhoto != "" {
				rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.UsersBulkName, userRes.LowQualityPhoto, minio.RemoveObjectOptions{})
				if rmLqErr != nil {
					return rmLqErr
				}
			}
			if userRes.Thumbnail != "" {
				rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.UsersBulkName, userRes.Thumbnail, minio.RemoveObjectOptions{})
				if rmThErr != nil {
					return rmThErr
				}
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(ctx, tx, &entity.User{ID: &userId}, &entity.User{HighQualityPhoto: req.User.HighQualityPhoto, LowQualityPhoto: req.User.LowQualityPhoto, Thumbnail: req.User.Thumbnail, BlurHash: req.User.BlurHash})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		} else if req.User.FullName != "" {
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(ctx, tx, &entity.User{ID: &userId}, &entity.User{FullName: req.User.FullName})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
	if updatedUserRes.HighQualityPhoto != "" {
		highQualityPhotoUrl = i.config.UsersBulkName + "/" + updatedUserRes.HighQualityPhoto
		lowQualityPhotoUrl = i.config.UsersBulkName + "/" + updatedUserRes.LowQualityPhoto
		thumbnailUrl = i.config.UsersBulkName + "/" + updatedUserRes.Thumbnail

	}
	return &pb.User{
		Id:                  updatedUserRes.ID.String(),
		FullName:            updatedUserRes.FullName,
		Email:               updatedUserRes.Email,
		HighQualityPhoto:    updatedUserRes.HighQualityPhoto,
		HighQualityPhotoUrl: highQualityPhotoUrl,
		LowQualityPhoto:     updatedUserRes.LowQualityPhoto,
		LowQualityPhotoUrl:  lowQualityPhotoUrl,
		Thumbnail:           updatedUserRes.Thumbnail,
		ThumbnailUrl:        thumbnailUrl,
		BlurHash:            updatedUserRes.BlurHash,
		UserAddress:         nil,
		Permissions:         nil,
		CreateTime:          timestamppb.New(updatedUserRes.CreateTime),
		UpdateTime:          timestamppb.New(updatedUserRes.UpdateTime),
	}, nil
}
