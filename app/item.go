package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	ut "github.com/daniarmas/api_go/utils"
	utils "github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *ItemServer) ListItem(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	var invalidBusinessId, invalidBusinessCollectionId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "BusinessId",
				Description: "The BusinessId field is not a valid uuid v4",
			}
		}
	}
	if req.BusinessCollectionId != "" {
		if !utils.IsValidUUID(&req.BusinessCollectionId) {
			invalidArgs = true
			invalidBusinessCollectionId = &epb.BadRequest_FieldViolation{
				Field:       "BusinessCollectionId",
				Description: "The BusinessCollectionId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		if invalidBusinessCollectionId != nil {
			st, _ = st.WithDetails(
				invalidBusinessCollectionId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.ListItem(ctx, req, md)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	var invalidId, invalidLocation *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Id == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "Id",
			Description: "The Id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "Id",
				Description: "The Id field is not a valid uuid v4",
			}
		}
	}
	if req.Location == nil {
		invalidArgs = true
		invalidLocation = &epb.BadRequest_FieldViolation{
			Field:       "Location",
			Description: "The Location field is required",
		}
	} else if req.Location != nil {
		if req.Location.Latitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "Location.Latitude",
				Description: "The Location.Latitude field is required",
			}
		} else if req.Location.Longitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "Location.Longitude",
				Description: "The Location.Longitude field is required",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.GetItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "record not found":
			st = status.New(codes.NotFound, "Item not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	var businessCollectionId uuid.UUID
	if req.BusinessCollectionId != "" {
		businessCollectionId = uuid.MustParse(req.BusinessCollectionId)
	}
	itemId := uuid.MustParse(req.Id)
	item, err := m.itemService.UpdateItem(&dto.UpdateItemRequest{ItemId: &itemId, Name: req.Name, Description: req.Description, Price: req.Price, HighQualityPhoto: req.HighQualityPhoto, HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash, LowQualityPhoto: req.LowQualityPhoto, LowQualityPhotoBlurHash: req.LowQualityPhotoBlurHash, Thumbnail: req.Thumbnail, ThumbnailBlurHash: req.ThumbnailBlurHash, Availability: req.Availability, Status: req.Status.String(), BusinessColletionId: &businessCollectionId, Metadata: &md})
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "business is open":
			st = status.New(codes.InvalidArgument, "Business is open")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "record not found":
			st = status.New(codes.NotFound, "Item not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.UpdateItemResponse{Item: &pb.Item{
		Id:                       item.ID.String(),
		Name:                     item.Name,
		Description:              item.Description,
		Price:                    item.Price,
		Availability:             int32(item.Availability),
		BusinessId:               item.BusinessId.String(),
		BusinessCollectionId:     item.BusinessCollectionId.String(),
		HighQualityPhoto:         item.HighQualityPhoto,
		HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
		LowQualityPhoto:          item.LowQualityPhoto,
		LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
		Thumbnail:                item.Thumbnail,
		ThumbnailBlurHash:        item.ThumbnailBlurHash,
		Cursor:                   item.Cursor,
		Status:                   *ut.ParseItemStatusType(&item.Status),
		CreateTime:               timestamppb.New(item.CreateTime),
		UpdateTime:               timestamppb.New(item.UpdateTime),
	}}, nil
}

func (m *ItemServer) SearchItem(ctx context.Context, req *pb.SearchItemRequest) (*pb.SearchItemResponse, error) {
	var st *status.Status
	response, err := m.itemService.SearchItem(req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), req.SearchMunicipalityType.String())
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	itemsResponse := make([]*pb.SearchItem, 0, len(*response.Items))
	for _, e := range *response.Items {
		itemsResponse = append(itemsResponse, &pb.SearchItem{
			Id:                e.ID.String(),
			Name:              e.Name,
			Thumbnail:         e.Thumbnail,
			ThumbnailBlurHash: e.ThumbnailBlurHash,
			Price:             e.Price,
			Cursor:            int32(e.Cursor),
			Status:            *ut.ParseItemStatusType(&e.Status),
		})
	}
	return &pb.SearchItemResponse{Items: itemsResponse, NextPage: response.NextPage, SearchMunicipalityType: *ut.ParseSearchMunicipalityType(response.SearchMunicipalityType)}, nil
}

func (m *ItemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	itemId := uuid.MustParse(req.Id)
	err := m.itemService.DeleteItem(&dto.DeleteItemRequest{ItemId: &itemId, Metadata: &md})
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "business is open":
			st = status.New(codes.InvalidArgument, "Business is open")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "item in the cart":
			st = status.New(codes.InvalidArgument, "Item in the cart")
		case "cartitem not found":
			st = status.New(codes.NotFound, "CartItem not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}

func (m *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	var st *status.Status
	businessId := uuid.MustParse(req.BusinessId)
	res, err := m.itemService.CreateItem(&dto.CreateItemRequest{Name: req.Name, Description: req.Description, Price: req.Price, BusinessCollectionId: req.BusinessCollectionId, HighQualityPhoto: req.HighQualityPhoto, HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash, LowQualityPhoto: req.LowQualityPhoto, LowQualityPhotoBlurHash: req.LowQualityPhotoBlurHash, Thumbnail: req.Thumbnail, ThumbnailBlurHash: req.ThumbnailBlurHash, BusinessId: &businessId, Metadata: &md})
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.CreateItemResponse{Item: &pb.Item{Id: res.ID.String(), Name: res.Name, Description: res.Description, Price: res.Price, Status: *ut.ParseItemStatusType(&res.Status), Availability: int32(res.Availability), BusinessId: res.BusinessId.String(), BusinessCollectionId: res.BusinessCollectionId.String(), HighQualityPhoto: res.HighQualityPhoto, HighQualityPhotoBlurHash: res.HighQualityPhotoBlurHash, LowQualityPhoto: res.LowQualityPhoto, LowQualityPhotoBlurHash: res.LowQualityPhotoBlurHash, Thumbnail: res.Thumbnail, ThumbnailBlurHash: res.ThumbnailBlurHash, CreateTime: timestamppb.New(res.CreateTime), UpdateTime: timestamppb.New(res.UpdateTime)}}, nil
}
