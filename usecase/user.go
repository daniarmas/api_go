package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, md *utils.ClientMetadata) (*pb.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, md *utils.ClientMetadata) (*pb.UpdateUserResponse, error)
	GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest, md *utils.ClientMetadata) (*pb.GetAddressInfoResponse, error)
}

type userService struct {
	dao    repository.DAO
	config *utils.Config
}

func NewUserService(dao repository.DAO, config *utils.Config) UserService {
	return &userService{dao: dao, config: config}
}

func (i *userService) GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest, md *utils.ClientMetadata) (*pb.GetAddressInfoResponse, error) {
	var res pb.GetAddressInfoResponse
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
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

func (i *userService) GetUser(ctx context.Context, md *utils.ClientMetadata) (*pb.GetUserResponse, error) {
	var userRes *models.User
	var userErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		userRes, userErr = i.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: authorizationTokenRes.UserId}, nil)
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
	return &pb.GetUserResponse{User: &pb.User{
		Id:                  userRes.ID.String(),
		FullName:            userRes.FullName,
		Email:               userRes.Email,
		HighQualityPhoto:    userRes.HighQualityPhoto,
		HighQualityPhotoUrl: i.config.UsersBulkName + "/" + userRes.HighQualityPhoto,
		LowQualityPhoto:     userRes.LowQualityPhoto,
		LowQualityPhotoUrl:  i.config.UsersBulkName + "/" + userRes.LowQualityPhoto,
		Thumbnail:           userRes.Thumbnail,
		ThumbnailUrl:        i.config.UsersBulkName + "/" + userRes.Thumbnail,
		BlurHash:            userRes.BlurHash,
		UserAddress:         userAddress,
		Permissions:         permissions,
		CreateTime:          timestamppb.New(userRes.CreateTime),
		UpdateTime:          timestamppb.New(userRes.UpdateTime),
	}}, nil
}

func (i *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, md *utils.ClientMetadata) (*pb.UpdateUserResponse, error) {
	var updatedUserRes *models.User
	var updatedUserErr error
	var userId uuid.UUID
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
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
		userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{"id", "high_quality_photo", "low_quality_photo", "thumbnail"})
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
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: userRes.Email, Code: req.Code, DeviceIdentifier: *md.DeviceIdentifier, Type: "ChangeUserEmail"}, &[]string{"id"})
			if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
				return errors.New("verification code not found")
			} else if verificationCodeErr != nil {
				return verificationCodeErr
			}
			_, err := i.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{ID: verificationCode.ID}, nil)
			if err != nil {
				return err
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: &userId}, &models.User{Email: req.User.Email})
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
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: &userId}, &models.User{HighQualityPhoto: req.User.HighQualityPhoto, LowQualityPhoto: req.User.LowQualityPhoto, Thumbnail: req.User.Thumbnail, BlurHash: req.User.BlurHash})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		} else if req.User.FullName != "" {
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: &userId}, &models.User{FullName: req.User.FullName})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{User: &pb.User{
		Id:                  updatedUserRes.ID.String(),
		FullName:            updatedUserRes.FullName,
		Email:               updatedUserRes.Email,
		HighQualityPhoto:    updatedUserRes.HighQualityPhoto,
		HighQualityPhotoUrl: i.config.UsersBulkName + "/" + updatedUserRes.HighQualityPhoto,
		LowQualityPhoto:     updatedUserRes.LowQualityPhoto,
		LowQualityPhotoUrl:  i.config.UsersBulkName + "/" + updatedUserRes.LowQualityPhoto,
		Thumbnail:           updatedUserRes.Thumbnail,
		ThumbnailUrl:        i.config.UsersBulkName + "/" + updatedUserRes.Thumbnail,
		BlurHash:            updatedUserRes.BlurHash,
		UserAddress:         nil,
		Permissions:         nil,
		CreateTime:          timestamppb.New(updatedUserRes.CreateTime),
		UpdateTime:          timestamppb.New(updatedUserRes.UpdateTime),
	}}, nil
}
