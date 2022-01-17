package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *BusinessServer) Feed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedResponse, error) {
	var st *status.Status
	business, err := m.businessService.Feed(&dto.FeedRequest{Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, ProvinceFk: req.ProvinceFk, MunicipalityFk: req.MunicipalityFk, NextPage: req.NextPage, SearchMunicipalityType: req.SearchMunicipalityType.String()})
	if err != nil {
		switch err.Error() {
		case "banned user":
			st = status.New(codes.PermissionDenied, "User banned")
		case "banned device":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		case "user already exists":
			st = status.New(codes.AlreadyExists, "User already exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
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

func (m *BusinessServer) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest) (*pb.GetBusinessResponse, error) {
	var st *status.Status
	business, err := m.businessService.GetBusiness(ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, req.Id)
	if err != nil {
		switch err.Error() {
		case "banned user":
			st = status.New(codes.PermissionDenied, "User banned")
		case "banned device":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "business not found":
			st = status.New(codes.NotFound, "Business not found")
		case "user already exists":
			st = status.New(codes.AlreadyExists, "User already exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	if err != nil {
		return nil, err
	}
	return &pb.GetBusinessResponse{Business: &pb.Business{Id: business.ID.String(), Name: business.Name, Description: business.Description, Address: business.Address, Phone: business.Phone, Email: business.Email, HighQualityPhoto: business.HighQualityPhoto, HighQualityPhotoBlurHash: business.HighQualityPhotoBlurHash, LowQualityPhoto: business.LowQualityPhoto, LowQualityPhotoBlurHash: business.LowQualityPhotoBlurHash, Thumbnail: business.Thumbnail, ThumbnailBlurHash: business.ThumbnailBlurHash, IsOpen: business.IsOpen, ToPickUp: business.ToPickUp, DeliveryPrice: float64(business.DeliveryPrice), HomeDelivery: business.HomeDelivery, ProvinceFk: business.ProvinceFk.String(), MunicipalityFk: business.MunicipalityFk.String(), BusinessBrandFk: business.BusinessBrandFk.String()}}, nil
}
