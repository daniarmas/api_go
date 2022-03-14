package usecase

import (
	// "context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessService interface {
	Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error)
	GetBusiness(request *dto.GetBusinessRequest) (*dto.GetBusinessResponse, error)
	CreateBusiness(request *dto.CreateBusinessRequest) (*dto.CreateBusinessResponse, error)
}

type businessService struct {
	dao repository.DAO
}

func NewBusinessService(dao repository.DAO) BusinessService {
	return &businessService{dao: dao}
}

func (i *businessService) CreateBusiness(request *dto.CreateBusinessRequest) (*dto.CreateBusinessResponse, error) {
	var businessRes *models.Business
	var businessErr error
	var response dto.CreateBusinessResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &models.BusinessUser{UserFk: authorizationTokenRes.UserFk}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		// _, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.HighQualityPhotoObject)
		// if hqErr != nil && hqErr.Error() == "ObjectMissing" {
		// 	return errors.New("HighQualityPhotoObject missing")
		// } else if hqErr != nil {
		// 	return hqErr
		// }
		// _, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.LowQualityPhotoObject)
		// if lqErr != nil && lqErr.Error() == "ObjectMissing" {
		// 	return errors.New("LowQualityPhotoObject missing")
		// } else if lqErr != nil {
		// 	return lqErr
		// }
		// _, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, request.ThumbnailObject)
		// if tnErr != nil && tnErr.Error() == "ObjectMissing" {
		// 	return errors.New("ThumbnailObject missing")
		// } else if tnErr != nil {
		// 	return tnErr
		// }
		businessRes, businessErr = i.dao.NewBusinessQuery().CreateBusiness(tx, &models.Business{
			Name:                     request.Name,
			Description:              request.Description,
			Address:                  request.Address,
			Phone:                    request.Phone,
			Email:                    request.Email,
			HighQualityPhoto:         datasource.Config.BusinessAvatarBulkName + "/" + request.HighQualityPhotoObject,
			HighQualityPhotoObject:   request.HighQualityPhotoObject,
			HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash,
			LowQualityPhoto:          datasource.Config.BusinessAvatarBulkName + "/" + request.LowQualityPhotoObject,
			LowQualityPhotoObject:    request.LowQualityPhotoObject,
			LowQualityPhotoBlurHash:  request.LowQualityPhotoBlurHash,
			Thumbnail:                datasource.Config.BusinessAvatarBulkName + "/" + request.ThumbnailObject,
			ThumbnailObject:          request.ThumbnailObject,
			ThumbnailBlurHash:        request.ThumbnailBlurHash,
			TimeMarginOrderMonth:     request.TimeMarginOrderMonth,
			TimeMarginOrderDay:       request.TimeMarginOrderDay,
			TimeMarginOrderHour:      request.TimeMarginOrderHour,
			TimeMarginOrderMinute:    request.TimeMarginOrderMinute,
			DeliveryPrice:            float32(request.DeliveryPrice),
			ToPickUp:                 request.ToPickUp,
			HomeDelivery:             request.HomeDelivery,
			Coordinates:              request.Coordinates,
			ProvinceFk:               uuid.MustParse(request.ProvinceFk),
			MunicipalityFk:           uuid.MustParse(request.MunicipalityFk),
			BusinessBrandFk:          uuid.MustParse(request.BusinessBrandFk),
		})
		if businessErr != nil {
			return businessErr
		}
		response.Business = businessRes
		var unionBusinessAndMunicipalities = make([]*models.UnionBusinessAndMunicipality, 0, len(request.Municipalities))
		for _, item := range request.Municipalities {
			unionBusinessAndMunicipalities = append(unionBusinessAndMunicipalities, &models.UnionBusinessAndMunicipality{
				BusinessFk:     businessRes.ID,
				MunicipalityFk: uuid.MustParse(item),
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
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, feedRequest.NextPage, false, feedRequest.HomeDelivery, feedRequest.ToPickUp)
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
				businessResAdd, businessErrAdd = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, int32(length), feedRequest.ProvinceFk, feedRequest.MunicipalityFk, 0, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
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
				businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, 0, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
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
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, feedRequest.NextPage, true, feedRequest.HomeDelivery, feedRequest.ToPickUp)
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
				Phone:                    e.Phone,
				Email:                    e.Email,
				HighQualityPhoto:         e.HighQualityPhoto,
				HighQualityPhotoBlurHash: e.HighQualityPhotoBlurHash,
				LowQualityPhoto:          e.LowQualityPhoto,
				LowQualityPhotoBlurHash:  e.LowQualityPhotoBlurHash,
				Thumbnail:                e.Thumbnail,
				ThumbnailBlurHash:        e.ThumbnailBlurHash,
				// IsOpen:                   e.IsOpen,
				DeliveryPrice:            e.DeliveryPrice,
				TimeMarginOrderMonth:     e.TimeMarginOrderMonth,
				TimeMarginOrderDay:       e.TimeMarginOrderDay,
				TimeMarginOrderHour:      e.TimeMarginOrderHour,
				TimeMarginOrderMinute:    e.TimeMarginOrderMinute,
				ToPickUp:                 e.ToPickUp,
				HomeDelivery:             e.HomeDelivery,
				BusinessBrandFk:          e.BusinessBrandFk,
				ProvinceFk:               e.ProvinceFk,
				Cursor:                   int32(e.Cursor),
				MunicipalityFk:           e.MunicipalityFk,
			})
		}
	}
	response.Businesses = &businessReturn
	return &response, nil
}

func (v *businessService) GetBusiness(request *dto.GetBusinessRequest) (*dto.GetBusinessResponse, error) {
	var businessRes *models.Business
	var itemCategoryRes *[]models.BusinessItemCategory
	var businessErr, itemCategoryErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		businessRes, businessErr = v.dao.NewBusinessQuery().GetBusiness(tx, &models.Business{Coordinates: request.Coordinates, ID: uuid.MustParse(request.Id)})
		if businessErr != nil {
			return businessErr
		}
		itemCategoryRes, itemCategoryErr = v.dao.NewItemCategoryQuery().ListItemCategory(tx, &models.BusinessItemCategory{BusinessFk: uuid.MustParse(request.Id)})
		if itemCategoryErr != nil {
			return itemCategoryErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.GetBusinessResponse{Business: businessRes, ItemCategory: itemCategoryRes}, nil
}
