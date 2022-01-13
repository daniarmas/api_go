package service

import (
	// "fmt"

	"github.com/daniarmas/api_go/datastruct"
	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
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
	var businessResAdd *[]datastruct.Business
	var businessErr, businessErrAdd error
	var response dto.FeedResponse
	if feedRequest.SearchMunicipalityType == pb.SearchMunicipalityType_More.String() {
		err := repository.DB.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, feedRequest.NextPage, false)
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
			err := repository.DB.Transaction(func(tx *gorm.DB) error {
				businessResAdd, businessErrAdd = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, int32(length), feedRequest.ProvinceFk, feedRequest.MunicipalityFk, 0, true)
				if businessErrAdd != nil {
					return businessErrAdd
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if len(*businessResAdd) > length {
				*businessResAdd = (*businessResAdd)[:len(*businessResAdd)-1]
			}
			*businessRes = append(*businessRes, *businessResAdd...)
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
		} else if len(*businessRes) == 0 {
			err := repository.DB.Transaction(func(tx *gorm.DB) error {
				businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, 0, true)
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
		err := repository.DB.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, feedRequest.Location, 5, feedRequest.ProvinceFk, feedRequest.MunicipalityFk, feedRequest.NextPage, true)
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
		} else if len(*businessRes) <= 5 && len(*businessRes) != 0 {
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
		}
		response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
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
	response.Businesses = &businessReturn
	return &response, nil
}
