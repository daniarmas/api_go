package usecase

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/repository"
)

type ObjectStorageService interface {
	GetPresignedPutObject(request *dto.GetPresignedPutObjectRequest) (*dto.GetPresignedPutObjectResponse, error)
}

type objectStorageService struct {
	dao repository.DAO
}

func NewObjectStorageService(dao repository.DAO) ObjectStorageService {
	return &objectStorageService{dao: dao}
}

func (i *objectStorageService) GetPresignedPutObject(request *dto.GetPresignedPutObjectRequest) (*dto.GetPresignedPutObjectResponse, error) {
	var bucket string
	var response dto.GetPresignedPutObjectResponse
	switch request.PhotoType {
	case "PhotoTypeBusiness":
		bucket = repository.Config.BusinessAvatarBulkName
	case "PhotoTypeItem":
		bucket = repository.Config.ItemsBulkName
	case "PhotoTypeUser":
		bucket = repository.Config.UsersBulkName
	}
	hqRes, hqErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, request.HighQualityPhotoObject, time.Duration(10)*time.Minute)
	if hqErr != nil {
		return nil, hqErr
	}
	lqRes, lqErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, request.LowQualityPhotoObject, time.Duration(10)*time.Minute)
	if lqErr != nil {
		return nil, lqErr
	}
	thRes, thErr := i.dao.NewObjectStorageRepository().PresignedPutObject(context.Background(), bucket, request.ThumbnailQualityPhotoObject, time.Duration(10)*time.Minute)
	if thErr != nil {
		return nil, thErr
	}
	response.HighQualityPhotoPresignedPutUrl = *hqRes
	response.LowQualityPhotoPresignedPutUrl = *lqRes
	response.ThumbnailQualityPhotoPresignedPutUrl = *thRes
	return &response, nil
}
