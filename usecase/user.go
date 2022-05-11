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
		provinceRes, provinceErr := i.dao.NewProvinceRepository().GetProvince(tx, &models.Province{ID: muncipalityRes.ProvinceFk})
		if provinceErr != nil && provinceErr.Error() == "record not found" {
			return errors.New("province not found")
		} else if provinceErr != nil {
			return provinceErr
		}
		response.MunicipalityId = muncipalityRes.ID
		response.MunicipalityName = muncipalityRes.Name
		response.ProvinceId = provinceRes.ID
		response.ProvinceName = provinceRes.Name
		response.ProvinceNameAbbreviation = provinceRes.ProvinceNameAbbreviation
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		user, userErr = i.dao.NewUserQuery().GetUserWithAddress(tx, &models.User{ID: *authorizationTokenRes.UserFk}, nil)
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserFk})
			if userErr != nil {
				return userErr
			}
			verificationCode, verificationCodeErr := i.dao.NewVerificationCodeQuery().GetVerificationCode(tx, &models.VerificationCode{Email: userRes.Email, Code: request.Code, DeviceId: request.Metadata.Get("deviceid")[0], Type: "ChangeUserEmail"}, &[]string{"id"})
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

		} else if request.Alias != "" {
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			user, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserFk})
			if userErr != nil && userErr.Error() == "record not found" {
				return errors.New("unauthenticated")
			} else if userErr != nil {
				return userErr
			}
			alias, aliasErr := i.dao.NewUserQuery().GetUser(tx, &models.User{Alias: request.Alias})
			if aliasErr != nil && aliasErr.Error() != "record not found" {
				return aliasErr
			} else if alias != nil && aliasErr == nil {
				return errors.New("user already exist")
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: user.ID}, &models.User{Alias: request.Alias})
			if updatedUserErr != nil {
				return updatedUserErr
			}

		} else if request.HighQualityPhotoObject != "" && request.HighQualityPhotoBlurHash != "" && request.LowQualityPhotoObject != "" && request.LowQualityPhotoBlurHash != "" && request.ThumbnailObject != "" && request.ThumbnailBlurHash != "" {
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserFk})
			if userErr != nil {
				return userErr
			}
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.HighQualityPhotoObject)
			if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.LowQualityPhotoObject)
			if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.UsersBulkName, request.ThumbnailObject)
			if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.HighQualityPhotoObject}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.HighQualityPhotoObject})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.LowQualityPhotoObject}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.LowQualityPhotoObject})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.UsersDeletedBulkName, Object: userRes.ThumbnailObject}, minio.CopySrcOptions{Bucket: repository.Config.UsersBulkName, Object: userRes.ThumbnailObject})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.HighQualityPhotoObject, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.LowQualityPhotoObject, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.UsersBulkName, userRes.ThumbnailObject, minio.RemoveObjectOptions{})
			if rmThErr != nil {
				return rmThErr
			}
			updatedUserRes, updatedUserErr = i.dao.NewUserQuery().UpdateUser(tx, &models.User{ID: userRes.ID}, &models.User{HighQualityPhotoObject: request.HighQualityPhotoObject, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, HighQualityPhoto: datasource.Config.UsersBulkName + "/" + request.HighQualityPhotoObject, LowQualityPhoto: datasource.Config.UsersBulkName + "/" + request.LowQualityPhotoObject, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, LowQualityPhotoObject: request.LowQualityPhotoObject, Thumbnail: datasource.Config.UsersBulkName + "/" + request.ThumbnailObject, ThumbnailObject: request.ThumbnailObject, ThumbnailBlurHash: request.ThumbnailBlurHash})
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
			authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if authorizationTokenRes == nil {
				return errors.New("unauthenticated")
			}
			userRes, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserFk})
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
