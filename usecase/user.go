package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/minio/minio-go/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, md *utils.ClientMetadata) (*pb.GetUserResponse, error)
	UpdateUser(request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	GetAddressInfo(request *dto.GetAddressInfoRequest) (*dto.GetAddressInfoResponse, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (i *userService) GetAddressInfo(request *dto.GetAddressInfoRequest) (*dto.GetAddressInfoResponse, error) {
	var response dto.GetAddressInfoResponse
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		muncipalityRes, err := i.dao.NewMunicipalityRepository().GetMunicipalityByCoordinate(tx, request.Coordinates)
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
		response.MunicipalityId = muncipalityRes.ID
		response.MunicipalityName = muncipalityRes.Name
		response.ProvinceId = provinceRes.ID
		response.ProvinceName = provinceRes.Name
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
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
			ResidenceType:  *utils.ParseResidenceType(item.ResidenceType),
			ProvinceId:     item.ProvinceId.String(),
			MunicipalityId: item.MunicipalityId.String(),
			Coordinates:    &pb.Point{Latitude: item.Coordinates.Coords()[0], Longitude: item.Coordinates.Coords()[1]},
			UserId:         item.UserId.String(),
			CreateTime:     timestamppb.New(item.CreateTime),
			UpdateTime:     timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.GetUserResponse{User: &pb.User{
		Id:                       userRes.ID.String(),
		FullName:                 userRes.FullName,
		Email:                    userRes.Email,
		PhoneNumber:              userRes.PhoneNumber,
		HighQualityPhoto:         userRes.HighQualityPhoto,
		HighQualityPhotoUrl:      userRes.HighQualityPhoto,
		HighQualityPhotoBlurHash: userRes.HighQualityPhotoBlurHash,
		LowQualityPhoto:          userRes.LowQualityPhoto,
		LowQualityPhotoUrl:       userRes.LowQualityPhoto,
		LowQualityPhotoBlurHash:  userRes.LowQualityPhotoBlurHash,
		Thumbnail:                userRes.Thumbnail,
		ThumbnailUrl:             userRes.Thumbnail,
		ThumbnailBlurHash:        userRes.ThumbnailBlurHash,
		UserAddress:              userAddress,
		Permissions:              permissions,
		CreateTime:               timestamppb.New(userRes.CreateTime),
		UpdateTime:               timestamppb.New(userRes.UpdateTime),
	}}, nil
}

func (i *userService) UpdateUser(request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	var updatedUserRes *models.User
	var updatedUserErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		if request.Email != "" && request.Code != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{})
			if userErr != nil {
				return userErr
			}
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: userRes.Email, Code: request.Code, DeviceIdentifier: request.Metadata.Get("deviceid")[0], Type: "ChangeUserEmail"}, &[]string{"id"})
			if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
				return errors.New("verification code not found")
			} else if verificationCodeErr != nil {
				return verificationCodeErr
			}
			_, err := i.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{ID: verificationCode.ID}, nil)
			if err != nil {
				return err
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: userRes.ID}, &models.User{Email: request.Email})
			if updatedUserErr != nil {
				return updatedUserErr
			}

		} else if request.HighQualityPhoto != "" && request.HighQualityPhotoBlurHash != "" && request.LowQualityPhoto != "" && request.LowQualityPhotoBlurHash != "" && request.Thumbnail != "" && request.ThumbnailBlurHash != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{})
			if userErr != nil {
				return userErr
			}
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.HighQualityPhoto)
			if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.LowQualityPhoto)
			if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.Thumbnail)
			if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.HighQualityPhoto})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.LowQualityPhoto})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.Thumbnail})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.Thumbnail, minio.RemoveObjectOptions{})
			if rmThErr != nil {
				return rmThErr
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: userRes.ID}, &models.User{HighQualityPhoto: request.HighQualityPhoto, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, LowQualityPhoto: request.LowQualityPhoto, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, Thumbnail: request.Thumbnail, ThumbnailBlurHash: request.ThumbnailBlurHash})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		} else if request.FullName != "" {
			jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserId}, &[]string{})
			if userErr != nil {
				return userErr
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: userRes.ID}, &models.User{FullName: request.FullName})
			if updatedUserErr != nil {
				return updatedUserErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.UpdateUserResponse{User: updatedUserRes}, nil
}
