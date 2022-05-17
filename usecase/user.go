package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(metadata *metadata.MD) (*models.User, error)
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
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		muncipalityRes, municipalityErr := i.dao.NewMunicipalityRepository().GetMunicipalityByCoordinate(tx, request.Coordinates)
		if municipalityErr != nil && municipalityErr.Error() == "record not found" {
			return errors.New("municipality not found")
		} else if municipalityErr != nil {
			return municipalityErr
		}
		provinceRes, provinceErr := i.dao.NewProvinceRepository().GetProvince(tx, &models.Province{ID: muncipalityRes.ProvinceId})
		if provinceErr != nil && provinceErr.Error() == "record not found" {
			return errors.New("province not found")
		} else if provinceErr != nil {
			return provinceErr
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

func (i *userService) GetUser(metadata *metadata.MD) (*models.User, error) {
	var user *models.User
	var userErr error
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		user, userErr = i.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: *authorizationTokenRes.UserId}, nil)
		if userErr != nil {
			return userErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *userService) UpdateUser(request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	var updatedUserRes *models.User
	var updatedUserErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserId}, &[]string{})
			if userErr != nil {
				return userErr
			}
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: userRes.Email, Code: request.Code, DeviceIdentifier: request.Metadata.Get("deviceid")[0], Type: "ChangeUserEmail"}, &[]string{"id"})
			if verificationCodeErr != nil && verificationCodeErr.Error() == "record not found" {
				return errors.New("verification code not found")
			} else if verificationCodeErr != nil {
				return verificationCodeErr
			}
			deleteVerificationCodeErr := i.dao.NewVerificationCodeQuery().DeleteVerificationCode(tx, &models.VerificationCode{ID: verificationCode.ID})
			if deleteVerificationCodeErr != nil {
				return deleteVerificationCodeErr
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
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserId}, &[]string{})
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
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserId}, &[]string{})
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
