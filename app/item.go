package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg/grpc"
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
				Field:       "businessId",
				Description: "The businessId field is not a valid uuid v4",
			}
		}
	}
	if req.BusinessCollectionId != "" {
		if !utils.IsValidUUID(&req.BusinessCollectionId) {
			invalidArgs = true
			invalidBusinessCollectionId = &epb.BadRequest_FieldViolation{
				Field:       "businessCollectionId",
				Description: "The businessCollectionId field is not a valid uuid v4",
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
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	var invalidId, invalidLocation *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Id == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "id",
			Description: "The id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "id",
				Description: "The id field is not a valid uuid v4",
			}
		}
	}
	if req.Location == nil {
		invalidArgs = true
		invalidLocation = &epb.BadRequest_FieldViolation{
			Field:       "location",
			Description: "The location field is required",
		}
	} else if req.Location != nil {
		if req.Location.Latitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "location.latitude",
				Description: "The location.latitude field is required",
			}
		} else if req.Location.Longitude == 0 {
			invalidArgs = true
			invalidLocation = &epb.BadRequest_FieldViolation{
				Field:       "location.longitude",
				Description: "The location.longitude field is required",
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "record not found":
			st = status.New(codes.NotFound, "Item not found")
		default:
			st = status.New(codes.Internal, err.Error())
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.Item, error) {
	var invalidMunicipalityId, invalidThumbnail, invalidBlurhash, invalidHighQualityPhoto, invalidLowQualityPhoto, invalidProvinceId, invalidBusinessCollectionId, invalidBusinessId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.Item.Thumbnail != "" || req.Item.HighQualityPhoto != "" || req.Item.LowQualityPhoto != "" {
		if req.Item.Thumbnail == "" {
			invalidArgs = true
			invalidThumbnail = &epb.BadRequest_FieldViolation{
				Field:       "item.thumbnail",
				Description: "The item.thumbnail field is required for update the item photo",
			}
		}
		if req.Item.HighQualityPhoto == "" {
			invalidArgs = true
			invalidHighQualityPhoto = &epb.BadRequest_FieldViolation{
				Field:       "item.highQualityPhoto",
				Description: "The item.highQualityPhoto field is required for update the item photo",
			}
		}
		if req.Item.LowQualityPhoto == "" {
			invalidArgs = true
			invalidLowQualityPhoto = &epb.BadRequest_FieldViolation{
				Field:       "item.lowQualityPhoto",
				Description: "The item.lowQualityPhoto field is required for update the item photo",
			}
		}
		if req.Item.BlurHash == "" {
			invalidArgs = true
			invalidBlurhash = &epb.BadRequest_FieldViolation{
				Field:       "item.blurHash",
				Description: "The item.blurHash field is required for update the item photo",
			}
		}
	}
	if req.Item.ProvinceId != "" {
		if !utils.IsValidUUID(&req.Item.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "item.provinceId",
				Description: "The item.provinceId field is not a valid uuid v4",
			}
		}
	}
	if req.Item.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.Item.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "item.municipalityId",
				Description: "The item.municipalityId field is not a valid uuid v4",
			}
		}
	}
	if req.Item.BusinessId != "" {
		if !utils.IsValidUUID(&req.Item.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "item.businessId",
				Description: "The item.businessId field is not a valid uuid v4",
			}
		}
	}
	if req.Item.BusinessCollectionId != "" {
		if !utils.IsValidUUID(&req.Item.BusinessCollectionId) {
			invalidArgs = true
			invalidBusinessCollectionId = &epb.BadRequest_FieldViolation{
				Field:       "item.businessCollectionId",
				Description: "The item.businessCollectionId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		if invalidProvinceId != nil {
			st, _ = st.WithDetails(
				invalidProvinceId,
			)
		}
		if invalidBusinessCollectionId != nil {
			st, _ = st.WithDetails(
				invalidBusinessCollectionId,
			)
		}
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		if invalidBlurhash != nil {
			st, _ = st.WithDetails(
				invalidBlurhash,
			)
		}
		if invalidLowQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidLowQualityPhoto,
			)
		}
		if invalidHighQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidHighQualityPhoto,
			)
		}
		if invalidThumbnail != nil {
			st, _ = st.WithDetails(
				invalidThumbnail,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.UpdateItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "item in the cart":
			st = status.New(codes.InvalidArgument, "Item in the cart")
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

func (m *ItemServer) SearchItem(ctx context.Context, req *pb.SearchItemRequest) (*pb.SearchItemResponse, error) {
	var invalidProvinceId, invalidMunicipalityId, invalidName, invalidSearchMunicipalityType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.SearchMunicipalityType == pb.SearchMunicipalityType_SearchMunicipalityTypeUnspecified {
		invalidArgs = true
		invalidSearchMunicipalityType = &epb.BadRequest_FieldViolation{
			Field:       "searchMunicipalityType",
			Description: "The searchMunicipalityType field is required",
		}
	}
	if req.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "name",
			Description: "The name field is required",
		}
	}
	if req.ProvinceId == "" {
		invalidArgs = true
		invalidProvinceId = &epb.BadRequest_FieldViolation{
			Field:       "provinceId",
			Description: "The provinceId field is required",
		}
	} else if req.ProvinceId != "" {
		if !utils.IsValidUUID(&req.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "provinceId",
				Description: "The provinceId field is not a valid uuid v4",
			}
		}
	}
	if req.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "municipalityId",
			Description: "The municipalityId field is required",
		}
	} else if req.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "municipalityId",
				Description: "The municipalityId field is not a valid uuid v4",
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
			Field:       "name",
			Description: "The name field is required",
		}
	}
	if req.BusinessId == "" {
		invalidArgs = true
		invalidBusinessId = &epb.BadRequest_FieldViolation{
			Field:       "businessId",
			Description: "The businessId field is required",
		}
	} else if req.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "businessId",
				Description: "The businessId field is not a valid uuid v4",
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil

}

func (m *ItemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*gp.Empty, error) {
	var invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.Id == "" {
		invalidArgs = true
		invalidId = &epb.BadRequest_FieldViolation{
			Field:       "id",
			Description: "The id field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "id",
				Description: "The id field is not a valid uuid v4",
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
		return nil, st.Err()
	}
	err := m.itemService.DeleteItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
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
		case "item not found":
			st = status.New(codes.NotFound, "Item not found")
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
	var invalidAvailability, invalidName, invalidPriceCup, invalidCostCup, invalidProfitCup, invalidPriceUsd, invalidCostUsd, invalidProfitUsd, invalidThumbnail, invalidHighQualityPhoto, invalidLowQualityPhoto, invalidBlurhash, invalidBusinessCollectionId, invalidBusinessId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.Item.Availability < -1 {
		invalidArgs = true
		invalidAvailability = &epb.BadRequest_FieldViolation{
			Field:       "item.availability",
			Description: "The item.availability field must be equal or greater than -1",
		}
	}
	if req.Item.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "item.name",
			Description: "The item.name field is required",
		}
	}
	if req.Item.PriceCup == "" {
		invalidArgs = true
		invalidThumbnail = &epb.BadRequest_FieldViolation{
			Field:       "item.priceCup",
			Description: "The item.priceCup field is required",
		}
	} else if req.Item.PriceCup != "" {
		if !utils.RegexpIsNumber(&req.Item.PriceCup) {
			invalidArgs = true
			invalidPriceCup = &epb.BadRequest_FieldViolation{
				Field:       "item.priceCup",
				Description: "The item.priceCup field is not a number",
			}
		}
	}
	if req.Item.CostCup != "" {
		if !utils.RegexpIsNumber(&req.Item.CostCup) {
			invalidArgs = true
			invalidCostCup = &epb.BadRequest_FieldViolation{
				Field:       "item.costCup",
				Description: "The item.costCup field is not a number",
			}
		}
	}
	if req.Item.ProfitCup != "" {
		if !utils.RegexpIsNumber(&req.Item.ProfitCup) {
			invalidArgs = true
			invalidProfitCup = &epb.BadRequest_FieldViolation{
				Field:       "item.profitCup",
				Description: "The item.profitCup field is not a number",
			}
		}
	}
	if req.Item.PriceUsd != "" {
		if !utils.RegexpIsNumber(&req.Item.PriceUsd) {
			invalidArgs = true
			invalidPriceUsd = &epb.BadRequest_FieldViolation{
				Field:       "item.priceUsd",
				Description: "The item.priceUsd field is not a number",
			}
		}
	}
	if req.Item.CostUsd != "" {
		if !utils.RegexpIsNumber(&req.Item.CostUsd) {
			invalidArgs = true
			invalidCostUsd = &epb.BadRequest_FieldViolation{
				Field:       "item.costUsd",
				Description: "The item.costUsd field is not a number",
			}
		}
	}
	if req.Item.ProfitUsd != "" {
		if !utils.RegexpIsNumber(&req.Item.ProfitUsd) {
			invalidArgs = true
			invalidProfitUsd = &epb.BadRequest_FieldViolation{
				Field:       "item.profitUsd",
				Description: "The item.profitUsd field is not a number",
			}
		}
	}
	if req.Item.Thumbnail == "" {
		invalidArgs = true
		invalidThumbnail = &epb.BadRequest_FieldViolation{
			Field:       "item.thumbnail",
			Description: "The item.thumbnail field is required",
		}
	}
	if req.Item.BlurHash == "" {
		invalidArgs = true
		invalidBlurhash = &epb.BadRequest_FieldViolation{
			Field:       "item.blurhash",
			Description: "The item.blurhash field is required",
		}
	}
	if req.Item.HighQualityPhoto == "" {
		invalidArgs = true
		invalidHighQualityPhoto = &epb.BadRequest_FieldViolation{
			Field:       "item.highQualityPhoto",
			Description: "The item.highQualityPhoto field is required",
		}
	}
	if req.Item.LowQualityPhoto == "" {
		invalidArgs = true
		invalidLowQualityPhoto = &epb.BadRequest_FieldViolation{
			Field:       "item.lowQualityPhoto",
			Description: "The item.lowQualityPhoto field is required",
		}
	}
	if req.Item.BusinessId == "" {
		invalidArgs = true
		invalidBusinessId = &epb.BadRequest_FieldViolation{
			Field:       "item.businessId",
			Description: "The item.businessId field is required",
		}
	} else if req.Item.BusinessId != "" {
		if !utils.IsValidUUID(&req.Item.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "item.businessId",
				Description: "The item.businessId field is not a valid uuid v4",
			}
		}
	}
	if req.Item.BusinessCollectionId != "" {
		if !utils.IsValidUUID(&req.Item.BusinessCollectionId) {
			invalidArgs = true
			invalidBusinessCollectionId = &epb.BadRequest_FieldViolation{
				Field:       "item.businessCollectionId",
				Description: "The item.businessCollectionId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidPriceCup != nil {
			st, _ = st.WithDetails(
				invalidPriceCup,
			)
		}
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		if invalidCostCup != nil {
			st, _ = st.WithDetails(
				invalidCostCup,
			)
		}
		if invalidPriceCup != nil {
			st, _ = st.WithDetails(
				invalidPriceCup,
			)
		}
		if invalidAvailability != nil {
			st, _ = st.WithDetails(
				invalidAvailability,
			)
		}
		if invalidProfitCup != nil {
			st, _ = st.WithDetails(
				invalidProfitCup,
			)
		}
		if invalidCostUsd != nil {
			st, _ = st.WithDetails(
				invalidCostUsd,
			)
		}
		if invalidPriceUsd != nil {
			st, _ = st.WithDetails(
				invalidPriceUsd,
			)
		}
		if invalidProfitUsd != nil {
			st, _ = st.WithDetails(
				invalidProfitUsd,
			)
		}
		if invalidBusinessCollectionId != nil {
			st, _ = st.WithDetails(
				invalidBusinessCollectionId,
			)
		}
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		if invalidLowQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidLowQualityPhoto,
			)
		}
		if invalidHighQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidHighQualityPhoto,
			)
		}
		if invalidThumbnail != nil {
			st, _ = st.WithDetails(
				invalidThumbnail,
			)
		}
		if invalidBlurhash != nil {
			st, _ = st.WithDetails(
				invalidBlurhash,
			)
		}
		return nil, st.Err()
	}
	res, err := m.itemService.CreateItem(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "highQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "lowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "thumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
