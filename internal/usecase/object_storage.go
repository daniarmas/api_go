package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/internal/repository"
	"github.com/daniarmas/api_go/utils"
	"gorm.io/gorm"
)

type ObjectStorageService interface {
	GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutObjectRequest, md *utils.ClientMetadata) (*pb.GetPresignedPutObjectResponse, error)
}

type objectStorageService struct {
	dao    repository.Repository
	config *config.Config
	sqldb  *sqldb.Sql
}

func NewObjectStorageService(dao repository.Repository, sqldb *sqldb.Sql, config *config.Config) ObjectStorageService {
	return &objectStorageService{dao: dao, sqldb: sqldb, config: config}
}

func (i *objectStorageService) GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutObjectRequest, md *utils.ClientMetadata) (*pb.GetPresignedPutObjectResponse, error) {
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
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
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var bucket string
	var response pb.GetPresignedPutObjectResponse
	switch req.PhotoType {
	case pb.PhotoType_PhotoTypeBusiness:
		bucket = i.config.BusinessAvatarBulkName
	case pb.PhotoType_PhotoTypeItem:
		bucket = i.config.ItemsBulkName
	case pb.PhotoType_PhotoTypeUser:
		bucket = i.config.UsersBulkName
	}
	hqRes, hqErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.HighQualityPhoto, time.Duration(10)*time.Minute)
	if hqErr != nil {
		return nil, hqErr
	}
	lqRes, lqErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.LowQualityPhoto, time.Duration(10)*time.Minute)
	if lqErr != nil {
		return nil, lqErr
	}
	thRes, thErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.ThumbnailQualityPhoto, time.Duration(10)*time.Minute)
	if thErr != nil {
		return nil, thErr
	}
	response.HighQualityPhotoPresignedPutUrl = *hqRes
	response.LowQualityPhotoPresignedPutUrl = *lqRes
	response.ThumbnailPresignedPutUrl = *thRes
	return &response, nil
}
