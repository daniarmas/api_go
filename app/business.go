package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

func (m *BusinessServer) Feed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedResponse, error) {
	business, err := m.businessService.Feed(&dto.FeedRequest{Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.Business, 0, len(*business.Businesses))
	for _, e := range *business.Businesses {
		itemsResponse = append(itemsResponse, &pb.Business{
			Id:                       e.ID.String(),
			Name:                     e.Name,
			Description:              e.Description,
			HighQualityPhoto:         e.HighQualityPhoto,
			HighQualityPhotoBlurHash: e.HighQualityPhotoBlurHash,
			LowQualityPhoto:          e.LowQualityPhoto,
			LowQualityPhotoBlurHash:  e.LowQualityPhotoBlurHash,
			Thumbnail:                e.Thumbnail,
			ThumbnailBlurHash:        e.ThumbnailBlurHash,
			Address:                  e.Address,
			Phone:                    e.Address,
			Email:                    e.Email,
			IsOpen:                   e.IsOpen,
			DeliveryPrice:            float64(e.DeliveryPrice),
			Coordinates:              &pb.Point{Latitude: e.Coordinates.Coords()[0], Longitude: e.Coordinates.Coords()[1]},
			Polygon:                  e.Polygon.FlatCoords(),
			LeadDayTime:              e.LeadDayTime,
			LeadHoursTime:            e.LeadHoursTime,
			LeadMinutesTime:          e.LeadMinutesTime,
			ToPickUp:                 e.ToPickUp,
			HomeDelivery:             e.HomeDelivery,
			BusinessBrandFk:          e.BusinessBrandFk.String(),
			ProvinceFk:               e.ProvinceFk.String(),
			MunicipalityFk:           e.MunicipalityFk.String(),
			Cursor:                   e.Cursor,
		})
	}
	return &pb.FeedResponse{Businesses: itemsResponse}, nil
}
