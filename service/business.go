package service

import (
	"github.com/daniarmas/api_go/datastruct"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/repository"
	"gorm.io/gorm"
)

type BusinessService interface {
	Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error)
}

type businessService struct {
	dao repository.DAO
}

func NewBusinessService(dao repository.DAO) BusinessService {
	return &businessService{dao: dao}
}

func (v *businessService) Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error) {
	var businessRes *[]datastruct.Business
	var businessErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		businessRes, businessErr = v.dao.NewBusinessQuery().ListBusiness(tx, &datastruct.Business{Coordinates: feedRequest.Location})
		if businessErr != nil {
			return businessErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	businessReturn := make([]dto.Business, 0, len(*businessRes))
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
			IsOpen:                   e.IsOpen,
			DeliveryPrice:            e.DeliveryPrice,
			Coordinates:              e.Coordinates,
			Polygon:                  e.Polygon,
			LeadDayTime:              e.LeadDayTime,
			LeadHoursTime:            e.LeadHoursTime,
			LeadMinutesTime:          e.LeadMinutesTime,
			ToPickUp:                 e.ToPickUp,
			HomeDelivery:             e.HomeDelivery,
			BusinessBrandFk:          e.BusinessBrandFk,
			ProvinceFk:               e.ProvinceFk,
			Cursor:                   int32(e.Cursor),
			MunicipalityFk:           e.MunicipalityFk,
		})
	}
	return &dto.FeedResponse{Businesses: &businessReturn}, nil
}
