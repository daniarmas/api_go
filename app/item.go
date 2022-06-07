package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
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

func (m *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.Item, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.itemService.UpdateItem(ctx, req, md)
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
	return res, nil
}

func (m *ItemServer) SearchItem(ctx context.Context, req *pb.SearchItemRequest) (*pb.SearchItemResponse, error) {
	var invalidProvinceId, invalidMunicipalityId, invalidName, invalidSearchMunicipalityType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.SearchMunicipalityType == pb.SearchMunicipalityType_SearchMunicipalityTypeUnspecified {
		invalidArgs = true
		invalidSearchMunicipalityType = &epb.BadRequest_FieldViolation{
			Field:       "SearchMunicipalityType",
			Description: "The SearchMunicipalityType field is required",
		}
	}
	if req.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "Name",
			Description: "The Name field is required",
		}
	}
	if req.ProvinceId == "" {
		invalidArgs = true
		invalidProvinceId = &epb.BadRequest_FieldViolation{
			Field:       "ProvinceId",
			Description: "The ProvinceId field is required",
		}
	} else if req.ProvinceId != "" {
		if !utils.IsValidUUID(&req.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "ProvinceId",
				Description: "The ProvinceId field is not a valid uuid v4",
			}
		}
	}
	if req.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "MunicipalityId",
			Description: "The MunicipalityId field is required",
		}
	} else if req.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "MunicipalityId",
				Description: "The MunicipalityId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidProvinceId != nil {
			st, _ = st.WithDetails(
				invalidProvinceId,
			)
		}
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		if invalidSearchMunicipalityType != nil {
			st, _ = st.WithDetails(
				invalidSearchMunicipalityType,
			)
		}
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.SearchItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *ItemServer) SearchItemByBusiness(ctx context.Context, req *pb.SearchItemByBusinessRequest) (*pb.SearchItemByBusinessResponse, error) {
	var invalidBusinessId, invalidName *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "Name",
			Description: "The Name field is required",
		}
	}
	if req.BusinessId == "" {
		invalidArgs = true
		invalidBusinessId = &epb.BadRequest_FieldViolation{
			Field:       "BusinessId",
			Description: "The BusinessId field is required",
		}
	} else if req.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "BusinessId",
				Description: "The BusinessId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.SearchItemByBusiness(ctx, req, md)
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil

}

func (m *ItemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*gp.Empty, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	err := m.itemService.DeleteItem(ctx, req, md)
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

func (m *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.Item, error) {
	md := utils.GetMetadata(ctx)
	var st *status.Status
	res, err := m.itemService.CreateItem(ctx, req, md)
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
	return res, nil
}
