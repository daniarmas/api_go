package usecase

import (
	"context"
	"time"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
)

type ObjectStorageService interface {
	GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutObjectRequest, md *utils.ClientMetadata) (*pb.GetPresignedPutObjectResponse, error)
}

type objectStorageService struct {
	dao repository.DAO
}

func NewObjectStorageService(dao repository.DAO) ObjectStorageService {
	return &objectStorageService{dao: dao}
}

func (i *objectStorageService) GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutObjectRequest, md *utils.ClientMetadata) (*pb.GetPresignedPutObjectResponse, error) {
	var bucket string
	var response pb.GetPresignedPutObjectResponse
	switch req.PhotoType {
	case pb.PhotoType_PhotoTypeBusiness:
		bucket = repository.Config.BusinessAvatarBulkName
	case pb.PhotoType_PhotoTypeItem:
		bucket = repository.Config.ItemsBulkName
	case pb.PhotoType_PhotoTypeUser:
		bucket = repository.Config.UsersBulkName
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
