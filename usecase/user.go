package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
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
	CreateUserAddress(ctx context.Context, req *pb.CreateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error)
	UpdateUserAddress(ctx context.Context, req *pb.UpdateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error)
	DeleteUserAddress(ctx context.Context, req *pb.DeleteUserAddressRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type userService struct {
	dao    repository.DAO
	config *utils.Config
	rdb    *redis.Client
}

func NewUserService(dao repository.DAO, config *utils.Config, rdb *redis.Client) UserService {
	return &userService{dao: dao, config: config, rdb: rdb}
}

func (i *userService) UpdateUserAddress(ctx context.Context, req *pb.UpdateUserAddressRequest, md *utils.ClientMetadata) (*pb.UserAddress, error) {
	var res pb.UserAddress
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr := i.dao.NewUserRepository().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{"id"})
		if userErr != nil {
			return userErr
		}
		listAddressRes, listAddressErr := i.dao.NewUserAddressRepository().ListUserAddress(tx, &models.UserAddress{UserId: userRes.ID}, &[]string{"id"})
		if listAddressErr != nil {
			return listAddressErr
		}
		if len(*listAddressRes) == 10 {
			return errors.New("only can have 10 user_address")
		}
		id := uuid.MustParse(req.Id)
		provinceId := uuid.MustParse(req.UserAddress.ProvinceId)
		municipalityId := uuid.MustParse(req.UserAddress.MunicipalityId)
		location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		updateUserAddressRes, updateUserAddressErr := i.dao.NewUserAddressRepository().UpdateUserAddress(tx, &models.UserAddress{ID: &id}, &models.UserAddress{Tag: req.UserAddress.Tag, Address: req.UserAddress.Address, Number: req.UserAddress.Number, Instructions: req.UserAddress.Instructions, UserId: userRes.ID, ProvinceId: &provinceId, MunicipalityId: &municipalityId, Coordinates: location})
		if updateUserAddressErr != nil && updateUserAddressErr.Error() == "record not found" {
			return errors.New("user address not found")
		} else if updateUserAddressErr != nil {
			return updateUserAddressErr
		}
		res = pb.UserAddress{
			Id:             updateUserAddressRes.ID.String(),
			Tag:            updateUserAddressRes.Tag,
			Number:         updateUserAddressRes.Number,
			Instructions:   updateUserAddressRes.Instructions,
			Address:        updateUserAddressRes.Address,
			UserId:         updateUserAddressRes.UserId.String(),
			ProvinceId:     updateUserAddressRes.ProvinceId.String(),
			MunicipalityId: updateUserAddressRes.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: updateUserAddressRes.Coordinates.FlatCoords()[0], Longitude: updateUserAddressRes.Coordinates.FlatCoords()[1]},
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		_, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		id := uuid.MustParse(req.Id)
		_, err = i.dao.NewUserAddressRepository().DeleteUserAddress(tx, &models.UserAddress{ID: &id}, nil)
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr := i.dao.NewUserRepository().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{"id"})
		if userErr != nil {
			return userErr
		}
		listAddressRes, listAddressErr := i.dao.NewUserAddressRepository().ListUserAddress(tx, &models.UserAddress{UserId: userRes.ID}, &[]string{"id"})
		if listAddressErr != nil {
			return listAddressErr
		}
		if len(*listAddressRes) == 10 {
			return errors.New("only can have 10 user_address")
		}
		provinceId := uuid.MustParse(req.UserAddress.ProvinceId)
		municipalityId := uuid.MustParse(req.UserAddress.MunicipalityId)
		location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.UserAddress.Coordinates.Latitude, req.UserAddress.Coordinates.Longitude}).SetSRID(4326)}
		createUserAddressRes, createUserAddressErr := i.dao.NewUserAddressRepository().CreateUserAddress(tx, &models.UserAddress{Tag: req.UserAddress.Tag, Address: req.UserAddress.Address, Number: req.UserAddress.Number, Instructions: req.UserAddress.Instructions, UserId: userRes.ID, ProvinceId: &provinceId, MunicipalityId: &municipalityId, Coordinates: location})
		if createUserAddressErr != nil {
			return createUserAddressErr
		}
		res = pb.UserAddress{
			Id:             createUserAddressRes.ID.String(),
			Tag:            createUserAddressRes.Tag,
			Number:         createUserAddressRes.Number,
			Instructions:   createUserAddressRes.Instructions,
			Address:        createUserAddressRes.Address,
			UserId:         createUserAddressRes.UserId.String(),
			ProvinceId:     createUserAddressRes.ProvinceId.String(),
			MunicipalityId: createUserAddressRes.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: createUserAddressRes.Coordinates.FlatCoords()[0], Longitude: createUserAddressRes.Coordinates.FlatCoords()[1]},
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userId := *authorizationTokenRes.UserId
		listUserAddressRes, listUserAddressErr := i.dao.NewUserAddressRepository().ListUserAddress(tx, &models.UserAddress{UserId: &userId}, nil)
		if listUserAddressErr != nil {
			return listUserAddressErr
		}
		userAddress := make([]*pb.UserAddress, 0, len(*listUserAddressRes))
		for _, i := range *listUserAddressRes {
			userAddress = append(userAddress, &pb.UserAddress{
				Id:             i.ID.String(),
				Tag:            i.Tag,
				UserId:         i.UserId.String(),
				Coordinates:    &pb.Point{Latitude: i.Coordinates.FlatCoords()[0], Longitude: i.Coordinates.FlatCoords()[1]},
				Address:        i.Address,
				Number:         i.Number,
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		muncipalityRes, err := i.dao.NewMunicipalityRepository().GetMunicipalityByCoordinate(tx, location)
		if err != nil && err.Error() == "record not found" {
			return errors.New("municipality not found")
		} else if err != nil {
			return err
		}
		provinceRes, err := i.dao.NewProvinceRepository().GetProvince(tx, &models.Province{ID: muncipalityRes.ProvinceId}, &[]string{})
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
	var userRes *models.User
	var userErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr = i.dao.NewUserRepository().GetUserWithAddress(tx, &models.User{ID: authorizationTokenRes.UserId}, nil)
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	userAddress := make([]*pb.UserAddress, 0, len(userRes.UserAddress))
	permissions := make([]*pb.Permission, 0, len(userRes.UserPermissions))
	for _, item := range userRes.UserPermissions {
		permissions = append(permissions, &pb.Permission{
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
	var updatedUserRes *models.User
	var updatedUserErr error
	var userId uuid.UUID
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		// chech if is the user or if have permission
		if req.User.Id != "" && authorizationTokenRes.UserId.String() != req.User.Id {
			_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "admin"}, &[]string{"id"})
			if err != nil && err.Error() == "record not found" {
				return errors.New("not have permission")
			} else if err != nil && err.Error() != "record not found" {
				return err
			}
		}
		userRes, userErr := i.dao.NewUserRepository().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{"id", "high_quality_photo", "low_quality_photo", "thumbnail"})
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
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeRepository().GetVerificationCode(tx, &models.VerificationCode{Email: userRes.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "ChangeUserEmail"}, &[]string{"id"})
			if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
				return errors.New("verification code not found")
			} else if verificationCodeErr != nil {
				return verificationCodeErr
			}
			_, err := i.dao.NewVerificationCodeRepository().DeleteVerificationCode(tx, &models.VerificationCode{ID: verificationCode.ID}, nil)
			if err != nil {
				return err
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(tx, &models.User{ID: &userId}, &models.User{Email: req.User.Email})
			if updatedUserErr != nil {
				return updatedUserErr
			}

		} else if req.User.HighQualityPhoto != "" && req.User.LowQualityPhoto != "" && req.User.Thumbnail != "" && req.User.BlurHash != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, req.User.HighQualityPhoto)
			if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, req.User.LowQualityPhoto)
			if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, req.User.Thumbnail)
			if tnErr != nil {
				return tnErr
			}
			if userRes.HighQualityPhoto != "" {
				_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.HighQualityPhoto})
				if copyHqErr != nil {
					return copyHqErr
				}
			}
			if userRes.LowQualityPhoto != "" {
				_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.LowQualityPhoto})
				if copyLqErr != nil {
					return copyLqErr
				}
			}
			if userRes.Thumbnail != "" {
				_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.Thumbnail})
				if copyThErr != nil {
					return copyThErr
				}
			}
			if userRes.HighQualityPhoto != "" {
				rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.HighQualityPhoto, minio.RemoveObjectOptions{})
				if rmHqErr != nil {
					return rmHqErr
				}
			}
			if userRes.LowQualityPhoto != "" {
				rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.LowQualityPhoto, minio.RemoveObjectOptions{})
				if rmLqErr != nil {
					return rmLqErr
				}
			}
			if userRes.Thumbnail != "" {
				rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.Thumbnail, minio.RemoveObjectOptions{})
				if rmThErr != nil {
					return rmThErr
				}
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(tx, &models.User{ID: &userId}, &models.User{HighQualityPhoto: req.User.HighQualityPhoto, LowQualityPhoto: req.User.LowQualityPhoto, Thumbnail: req.User.Thumbnail, BlurHash: req.User.BlurHash})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		} else if req.User.FullName != "" {
			updatedUserRes, updatedUserErr = i.dao.NewUserRepository().UpdateUser(tx, &models.User{ID: &userId}, &models.User{FullName: req.User.FullName})
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
