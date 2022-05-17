package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type BusinessService interface {
	Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error)
	GetBusiness(request *dto.GetBusinessRequest) (*dto.GetBusinessResponse, error)
	CreateBusiness(request *dto.CreateBusinessRequest) (*dto.CreateBusinessResponse, error)
	UpdateBusiness(request *dto.UpdateBusinessRequest) (*models.Business, error)
}

type businessService struct {
	dao repository.DAO
}

func NewBusinessService(dao repository.DAO) BusinessService {
	return &businessService{dao: dao}
}

func (i *businessService) UpdateBusiness(request *dto.UpdateBusinessRequest) (*models.Business, error) {
	var businessRes *models.Business
	var businessErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &models.BusinessUser{UserId: *authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		businessIsOpenRes, businessIsOpenErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: request.Id}, "OrderTypePickUp")
		if businessIsOpenErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessIsOpenErr
		} else if businessIsOpenRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: request.Id}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{BusinessId: request.Id})
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		getBusinessRes, getBusinessErr := i.dao.NewBusinessQuery().GetBusiness(tx, &models.Business{ID: request.Id})
		if getBusinessErr != nil {
			return getBusinessErr
		}
		if request.HighQualityPhoto != "" || request.LowQualityPhoto != "" || request.Thumbnail != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.HighQualityPhoto)
			if hqErr != nil && hqErr.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.LowQualityPhoto)
			if lqErr != nil && lqErr.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.Thumbnail)
			if tnErr != nil && tnErr.Error() == "ObjectMissing" {
				return errors.New("ThumbnailObject missing")
			} else if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.HighQualityPhoto})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.LowQualityPhoto})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.Thumbnail})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.Thumbnail, minio.RemoveObjectOptions{})
			if rmThErr != nil {
				return rmThErr
			}
		}
		var provinceId uuid.UUID
		var municipalityId uuid.UUID
		if request.ProvinceId != "" {
			provinceId = uuid.MustParse(request.ProvinceId)
		}
		if request.MunicipalityId != "" {
			municipalityId = uuid.MustParse(request.MunicipalityId)
		}
		businessRes, businessErr = i.dao.NewBusinessQuery().UpdateBusiness(tx, &models.Business{
			Name:                     request.Name,
			Description:              request.Description,
			Address:                  request.Address,
			HighQualityPhoto:         datasource.Config.BusinessAvatarBulkName + "/" + request.HighQualityPhoto,
			HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash,
			LowQualityPhoto:          datasource.Config.BusinessAvatarBulkName + "/" + request.LowQualityPhoto,
			LowQualityPhotoBlurHash:  request.LowQualityPhotoBlurHash,
			Thumbnail:                datasource.Config.BusinessAvatarBulkName + "/" + request.Thumbnail,
			ThumbnailBlurHash:        request.ThumbnailBlurHash,
			TimeMarginOrderMonth:     request.TimeMarginOrderMonth,
			TimeMarginOrderDay:       request.TimeMarginOrderDay,
			TimeMarginOrderHour:      request.TimeMarginOrderHour,
			TimeMarginOrderMinute:    request.TimeMarginOrderMinute,
			DeliveryPrice:            request.DeliveryPrice,
			ToPickUp:                 request.ToPickUp,
			HomeDelivery:             request.HomeDelivery,
			ProvinceId:               provinceId,
			MunicipalityId:           municipalityId,
		}, &models.Business{ID: request.Id})
		if businessErr != nil {
			return businessErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return businessRes, nil
}

func (i *businessService) CreateBusiness(request *dto.CreateBusinessRequest) (*dto.CreateBusinessResponse, error) {
	var businessRes *models.Business
	var businessErr error
	var response dto.CreateBusinessResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &models.BusinessUser{UserId: *authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("HighQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("LowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("ThumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		businessRes, businessErr = i.dao.NewBusinessQuery().CreateBusiness(tx, &models.Business{
			Name:                     request.Name,
			Description:              request.Description,
			Address:                  request.Address,
			HighQualityPhoto:         request.HighQualityPhoto,
			HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash,
			LowQualityPhoto:          request.LowQualityPhoto,
			LowQualityPhotoBlurHash:  request.LowQualityPhotoBlurHash,
			Thumbnail:                request.Thumbnail,
			ThumbnailBlurHash:        request.ThumbnailBlurHash,
			TimeMarginOrderMonth:     request.TimeMarginOrderMonth,
			TimeMarginOrderDay:       request.TimeMarginOrderDay,
			TimeMarginOrderHour:      request.TimeMarginOrderHour,
			TimeMarginOrderMinute:    request.TimeMarginOrderMinute,
			DeliveryPrice:            request.DeliveryPrice,
			ToPickUp:                 request.ToPickUp,
			HomeDelivery:             request.HomeDelivery,
			Coordinates:              request.Coordinates,
			ProvinceId:               uuid.MustParse(request.ProvinceId),
			MunicipalityId:           uuid.MustParse(request.MunicipalityId),
			BusinessBrandId:          uuid.MustParse(request.BusinessBrandId),
		})
		if businessErr != nil {
			return businessErr
		}
		response.Business = businessRes
		var unionBusinessAndMunicipalities = make([]*models.UnionBusinessAndMunicipality, 0, len(request.Municipalities))
		for _, item := range request.Municipalities {
			unionBusinessAndMunicipalities = append(unionBusinessAndMunicipalities, &models.UnionBusinessAndMunicipality{
				BusinessId:     businessRes.ID,
				MunicipalityId: uuid.MustParse(item),
			})
		}
		unionBusinessAndMunicipalityRes, unionBusinessAndMunicipalityErr := i.dao.NewUnionBusinessAndMunicipalityRepository().BatchCreateUnionBusinessAndMunicipality(tx, unionBusinessAndMunicipalities)
		if unionBusinessAndMunicipalityErr != nil {
			return unionBusinessAndMunicipalityErr
		}
		unionBusinessAndMunicipalityIds := make([]string, 0, len(unionBusinessAndMunicipalityRes))
		for _, item := range unionBusinessAndMunicipalityRes {
			unionBusinessAndMunicipalityIds = append(unionBusinessAndMunicipalityIds, item.ID.String())
		}
		unionBusinessAndMunicipalityWithMunicipalityRes, unionBusinessAndMunicipalityWithMunicipalityErr := i.dao.NewUnionBusinessAndMunicipalityRepository().ListUnionBusinessAndMunicipalityWithMunicipality(tx, unionBusinessAndMunicipalityIds)
		if unionBusinessAndMunicipalityWithMunicipalityErr != nil {
			return unionBusinessAndMunicipalityWithMunicipalityErr
		}
		response.UnionBusinessAndMunicipalityWithMunicipality = unionBusinessAndMunicipalityWithMunicipalityRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *businessService) Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error) {
	var businessRes *[]models.Business
	var businessResAdd *[]models.Business
	var businessErr, businessErrAdd error
	var response dto.FeedResponse
	if feedRequest.SearchMunicipalityType == pb.SearchMunicipalityType_More.String() {
		err := datasource.DB.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceId, feedRequest.MunicipalityId, feedRequest.NextPage, false, feedRequest.HomeDelivery, feedRequest.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if len(*businessRes) > 5 {
			*businessRes = (*businessRes)[:len(*businessRes)-1]
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			response.SearchMunicipalityType = pb.SearchMunicipalityType_More.String()
		} else if len(*businessRes) <= 5 && len(*businessRes) != 0 {
			length := 5 - len(*businessRes)
			err := datasource.DB.Transaction(func(tx *gorm.DB) error {
				businessResAdd, businessErrAdd = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, int32(length), feedRequest.ProvinceId, feedRequest.MunicipalityId, 0, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
				if businessErrAdd != nil {
					return businessErrAdd
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if businessResAdd != nil {
				if len(*businessResAdd) > length {
					*businessResAdd = (*businessResAdd)[:len(*businessResAdd)-1]
				}
				*businessRes = append(*businessRes, *businessResAdd...)
			}
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
		} else if len(*businessRes) == 0 {
			err := datasource.DB.Transaction(func(tx *gorm.DB) error {
				businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceId, feedRequest.MunicipalityId, 0, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
				if businessErr != nil {
					return businessErr
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if len(*businessRes) > 5 {
				*businessRes = (*businessRes)[:len(*businessRes)-1]
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
				response.SearchMunicipalityType = pb.SearchMunicipalityType_More.String()
			}
		}
	} else {
		err := datasource.DB.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceId, feedRequest.MunicipalityId, feedRequest.NextPage, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if businessRes != nil && len(*businessRes) > 5 {
			*businessRes = (*businessRes)[:len(*businessRes)-1]
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
		} else if businessRes != nil && len(*businessRes) <= 5 && len(*businessRes) != 0 {
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
		} else {
			response.NextPage = feedRequest.NextPage
		}
		response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
	}
	var businessReturn []dto.Business
	if businessRes != nil {
		businessReturn = make([]dto.Business, 0, len(*businessRes))
		for _, e := range *businessRes {
			businessReturn = append(businessReturn, dto.Business{
				ID:                       e.ID,
				Name:                     e.Name,
				Description:              e.Description,
				Address:                  e.Address,
				HighQualityPhoto:         e.HighQualityPhoto,
				HighQualityPhotoBlurHash: e.HighQualityPhotoBlurHash,
				LowQualityPhoto:          e.LowQualityPhoto,
				LowQualityPhotoBlurHash:  e.LowQualityPhotoBlurHash,
				Thumbnail:                e.Thumbnail,
				ThumbnailBlurHash:        e.ThumbnailBlurHash,
				DeliveryPrice:            e.DeliveryPrice,
				TimeMarginOrderMonth:     e.TimeMarginOrderMonth,
				TimeMarginOrderDay:       e.TimeMarginOrderDay,
				TimeMarginOrderHour:      e.TimeMarginOrderHour,
				TimeMarginOrderMinute:    e.TimeMarginOrderMinute,
				ToPickUp:                 e.ToPickUp,
				HomeDelivery:             e.HomeDelivery,
				BusinessBrandId:          e.BusinessBrandId,
				ProvinceId:               e.ProvinceId,
				Cursor:                   int32(e.Cursor),
				MunicipalityId:           e.MunicipalityId,
			})
		}
	}
	response.Businesses = &businessReturn
	return &response, nil
}

func (v *businessService) GetBusiness(request *dto.GetBusinessRequest) (*dto.GetBusinessResponse, error) {
	var businessRes *models.Business
	var businessCollectionRes *[]models.BusinessCollection
	var businessErr, itemCategoryErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		businessRes, businessErr = v.dao.NewBusinessQuery().GetBusiness(tx, &models.Business{ID: uuid.MustParse(request.Id)})
		if businessErr != nil {
			return businessErr
		}
		businessCollectionRes, itemCategoryErr = v.dao.NewBusinessCollectionQuery().ListBusinessCollection(tx, &models.BusinessCollection{BusinessId: uuid.MustParse(request.Id)})
		if itemCategoryErr != nil {
			return itemCategoryErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.GetBusinessResponse{Business: businessRes, BusinessCollections: businessCollectionRes}, nil
}
