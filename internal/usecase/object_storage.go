package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
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
	var res pb.GetPresignedPutObjectResponse
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
		var bucket string
		switch req.PhotoType {
		case pb.PhotoType_PhotoTypeBusiness:
			bucket = i.config.BusinessAvatarBulkName
		case pb.PhotoType_PhotoTypeItem:
			bucket = i.config.ItemsBulkName
		case pb.PhotoType_PhotoTypeUser:
			bucket = i.config.UsersBulkName
		}
		hqRes, err := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.HighQualityPhoto, time.Duration(10)*time.Minute)
		if err != nil {
			return err
		}
		lqRes, err := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.LowQualityPhoto, time.Duration(10)*time.Minute)
		if err != nil {
			return err
		}
		thRes, err := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, req.ThumbnailQualityPhoto, time.Duration(10)*time.Minute)
		if err != nil {
			return err
		}
		res.HighQualityPhotoPresignedPutUrl = *hqRes
		res.LowQualityPhotoPresignedPutUrl = *lqRes
		res.ThumbnailPresignedPutUrl = *thRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
