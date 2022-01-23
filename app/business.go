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
	getBusiness, err := m.businessService.GetBusiness(ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, req.Id)
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
	itemsCategoryResponse := make([]*pb.ItemCategory, 0, len(*getBusiness.ItemCategory))
	for _, e := range *getBusiness.ItemCategory {
		itemsCategoryResponse = append(itemsCategoryResponse, &pb.ItemCategory{
			Id:         e.ID.String(),
			Name:       e.Name,
			BusinessFk: e.BusinessFk.String(),
			Index:      e.Index,
			CreateTime: e.CreateTime.String(),
			UpdateTime: e.UpdateTime.String(),
		})
	}
	return &pb.GetBusinessResponse{Business: &pb.Business{Id: getBusiness.Business.ID.String(), Name: getBusiness.Business.Name, Description: getBusiness.Business.Description, Address: getBusiness.Business.Address, Phone: getBusiness.Business.Phone, Email: getBusiness.Business.Email, HighQualityPhoto: getBusiness.Business.HighQualityPhoto, HighQualityPhotoBlurHash: getBusiness.Business.HighQualityPhotoBlurHash, LowQualityPhoto: getBusiness.Business.LowQualityPhoto, LowQualityPhotoBlurHash: getBusiness.Business.LowQualityPhotoBlurHash, Thumbnail: getBusiness.Business.Thumbnail, ThumbnailBlurHash: getBusiness.Business.ThumbnailBlurHash, IsOpen: getBusiness.Business.IsOpen, ToPickUp: getBusiness.Business.ToPickUp, DeliveryPrice: float64(getBusiness.Business.DeliveryPrice), HomeDelivery: getBusiness.Business.HomeDelivery, ProvinceFk: getBusiness.Business.ProvinceFk.String(), MunicipalityFk: getBusiness.Business.MunicipalityFk.String(), BusinessBrandFk: getBusiness.Business.BusinessBrandFk.String()}, ItemCategory: itemsCategoryResponse}, nil
}
