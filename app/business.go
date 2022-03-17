package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *BusinessServer) CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest) (*pb.CreateBusinessResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.businessService.CreateBusiness(&dto.CreateBusinessRequest{
		Name:                     req.Name,
		Description:              req.Description,
		Address:                  req.Address,
		Phone:                    req.Phone,
		Email:                    req.Email,
		HighQualityPhotoObject:   req.HighQualityPhotoObject,
		HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash,
		LowQualityPhotoObject:    req.LowQualityPhotoObject,
		LowQualityPhotoBlurHash:  req.LowQualityPhotoBlurHash,
		ThumbnailObject:          req.ThumbnailObject,
		ThumbnailBlurHash:        req.HighQualityPhotoBlurHash,
		DeliveryPrice:            req.DeliveryPrice,
		TimeMarginOrderMonth:     req.TimeMarginOrderMonth,
		TimeMarginOrderDay:       req.TimeMarginOrderDay,
		TimeMarginOrderHour:      req.TimeMarginOrderHour,
		TimeMarginOrderMinute:    req.TimeMarginOrderMinute,
		Coordinates:              ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)},
		ToPickUp:                 req.ToPickUp,
		HomeDelivery:             req.HomeDelivery,
		BusinessBrandFk:          req.BusinessBrandFk,
		ProvinceFk:               req.ProvinceFk,
		MunicipalityFk:           req.MunicipalityFk,
		Metadata:                 &md,
		Municipalities:           req.Municipalities,
	})
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	municipalities := make([]*pb.UnionBusinessAndMunicipality, 0, len(*res.UnionBusinessAndMunicipalityWithMunicipality))
	for _, item := range *res.UnionBusinessAndMunicipalityWithMunicipality {
		municipalities = append(municipalities, &pb.UnionBusinessAndMunicipality{
			Id:             item.ID.String(),
			Name:           item.MunicipalityName,
			MunicipalityFk: item.MunicipalityFk.String(),
		})
	}
	return &pb.CreateBusinessResponse{Business: &pb.Business{Id: res.Business.ID.String(), Name: res.Business.Name, Description: res.Business.Description, Address: res.Business.Address, Phone: res.Business.Phone, Email: res.Business.Email, HighQualityPhoto: res.Business.HighQualityPhoto, HighQualityPhotoBlurHash: res.Business.HighQualityPhotoBlurHash, LowQualityPhoto: res.Business.LowQualityPhoto, LowQualityPhotoBlurHash: res.Business.LowQualityPhotoBlurHash, Thumbnail: res.Business.Thumbnail, ThumbnailBlurHash: res.Business.ThumbnailBlurHash, DeliveryPrice: float64(res.Business.DeliveryPrice), TimeMarginOrderMonth: res.Business.TimeMarginOrderMonth, TimeMarginOrderDay: res.Business.TimeMarginOrderDay, TimeMarginOrderHour: res.Business.TimeMarginOrderHour, TimeMarginOrderMinute: res.Business.TimeMarginOrderMinute, ToPickUp: res.Business.ToPickUp, HomeDelivery: res.Business.HomeDelivery, ProvinceFk: res.Business.ProvinceFk.String(), MunicipalityFk: res.Business.MunicipalityFk.String(), BusinessBrandFk: res.Business.BusinessBrandFk.String(), CreateTime: timestamppb.New(res.Business.CreateTime), UpdateTime: timestamppb.New(res.Business.UpdateTime), Coordinates: &pb.Point{Latitude: res.Business.Coordinates.FlatCoords()[0], Longitude: res.Business.Coordinates.FlatCoords()[1]}}, Municipalities: municipalities}, nil
}

func (m *BusinessServer) UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest) (*pb.UpdateBusinessResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.businessService.UpdateBusiness(&dto.UpdateBusinessRequest{
		Id:                       uuid.MustParse(req.Id),
		Name:                     req.Name,
		Description:              req.Description,
		Address:                  req.Address,
		Phone:                    req.Phone,
		Email:                    req.Email,
		HighQualityPhotoObject:   req.HighQualityPhotoObject,
		HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash,
		LowQualityPhotoObject:    req.LowQualityPhotoObject,
		LowQualityPhotoBlurHash:  req.LowQualityPhotoBlurHash,
		ThumbnailObject:          req.ThumbnailObject,
		ThumbnailBlurHash:        req.HighQualityPhotoBlurHash,
		DeliveryPrice:            req.DeliveryPrice,
		TimeMarginOrderMonth:     req.TimeMarginOrderMonth,
		TimeMarginOrderDay:       req.TimeMarginOrderDay,
		TimeMarginOrderHour:      req.TimeMarginOrderHour,
		TimeMarginOrderMinute:    req.TimeMarginOrderMinute,
		ToPickUp:                 req.ToPickUp,
		HomeDelivery:             req.HomeDelivery,
		ProvinceFk:               req.ProvinceFk,
		MunicipalityFk:           req.MunicipalityFk,
		Metadata:                 &md,
	})
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "business is open":
			st = status.New(codes.InvalidArgument, "Business is open")
		case "item in the cart":
			st = status.New(codes.InvalidArgument, "Item in the cart")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.UpdateBusinessResponse{Business: &pb.Business{Id: res.ID.String(), Name: res.Name, Description: res.Description, Address: res.Address, Phone: res.Phone, Email: res.Email, HighQualityPhoto: res.HighQualityPhoto, HighQualityPhotoBlurHash: res.HighQualityPhotoBlurHash, LowQualityPhoto: res.LowQualityPhoto, LowQualityPhotoBlurHash: res.LowQualityPhotoBlurHash, Thumbnail: res.Thumbnail, ThumbnailBlurHash: res.ThumbnailBlurHash, DeliveryPrice: float64(res.DeliveryPrice), TimeMarginOrderMonth: res.TimeMarginOrderMonth, TimeMarginOrderDay: res.TimeMarginOrderDay, TimeMarginOrderHour: res.TimeMarginOrderHour, TimeMarginOrderMinute: res.TimeMarginOrderMinute, ToPickUp: res.ToPickUp, HomeDelivery: res.HomeDelivery, ProvinceFk: res.ProvinceFk.String(), MunicipalityFk: res.MunicipalityFk.String(), BusinessBrandFk: res.BusinessBrandFk.String(), CreateTime: timestamppb.New(res.CreateTime), UpdateTime: timestamppb.New(res.UpdateTime)}}, nil
}

func (m *BusinessServer) Feed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	feedBusiness, err := m.businessService.Feed(&dto.FeedRequest{Location: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, ProvinceFk: req.ProvinceFk, MunicipalityFk: req.MunicipalityFk, HomeDelivery: req.HomeDelivery, ToPickUp: req.ToPickUp, NextPage: req.NextPage, SearchMunicipalityType: req.SearchMunicipalityType.String(), Metadata: &md})
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
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
	itemsResponse := make([]*pb.Business, 0, len(*feedBusiness.Businesses))
	for _, e := range *feedBusiness.Businesses {
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
			TimeMarginOrderMonth:     e.TimeMarginOrderMonth,
			TimeMarginOrderDay:       e.TimeMarginOrderDay,
			TimeMarginOrderHour:      e.TimeMarginOrderHour,
			TimeMarginOrderMinute:    e.TimeMarginOrderMinute,
			ToPickUp:                 e.ToPickUp,
			HomeDelivery:             e.HomeDelivery,
			BusinessBrandFk:          e.BusinessBrandFk.String(),
			ProvinceFk:               e.ProvinceFk.String(),
			MunicipalityFk:           e.MunicipalityFk.String(),
			Cursor:                   e.Cursor,
		})
	}
	return &pb.FeedResponse{Businesses: itemsResponse, SearchMunicipalityType: *utils.ParseSearchMunicipalityType(feedBusiness.SearchMunicipalityType), NextPage: feedBusiness.NextPage}, nil
}

func (m *BusinessServer) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest) (*pb.GetBusinessResponse, error) {
	var st *status.Status
	getBusiness, err := m.businessService.GetBusiness(&dto.GetBusinessRequest{Id: req.Id, Coordinates: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}})
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
			CreateTime: timestamppb.New(e.CreateTime),
			UpdateTime: timestamppb.New(e.UpdateTime),
		})
	}
	return &pb.GetBusinessResponse{Business: &pb.Business{Id: getBusiness.Business.ID.String(), Name: getBusiness.Business.Name, Description: getBusiness.Business.Description, Address: getBusiness.Business.Address, Phone: getBusiness.Business.Phone, Email: getBusiness.Business.Email, HighQualityPhoto: getBusiness.Business.HighQualityPhoto, HighQualityPhotoBlurHash: getBusiness.Business.HighQualityPhotoBlurHash, LowQualityPhoto: getBusiness.Business.LowQualityPhoto, LowQualityPhotoBlurHash: getBusiness.Business.LowQualityPhotoBlurHash, Thumbnail: getBusiness.Business.Thumbnail, ThumbnailBlurHash: getBusiness.Business.ThumbnailBlurHash, ToPickUp: getBusiness.Business.ToPickUp, DeliveryPrice: float64(getBusiness.Business.DeliveryPrice), HomeDelivery: getBusiness.Business.HomeDelivery, ProvinceFk: getBusiness.Business.ProvinceFk.String(), MunicipalityFk: getBusiness.Business.MunicipalityFk.String(), BusinessBrandFk: getBusiness.Business.BusinessBrandFk.String(), IsInRange: getBusiness.Business.IsInRange, Coordinates: &pb.Point{Latitude: getBusiness.Business.Coordinates.Coords()[1], Longitude: getBusiness.Business.Coordinates.Coords()[0]}, Distance: getBusiness.Business.Distance}, ItemCategory: itemsCategoryResponse}, nil
}
