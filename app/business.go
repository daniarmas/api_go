package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *BusinessServer) CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest) (*pb.CreateBusinessResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.businessService.CreateBusiness(ctx, req, md)
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
	// municipalities := make([]*pb.UnionBusinessAndMunicipality, 0, len(*res.UnionBusinessAndMunicipalityWithMunicipality))
	// for _, item := range *res.UnionBusinessAndMunicipalityWithMunicipality {
	// 	municipalities = append(municipalities, &pb.UnionBusinessAndMunicipality{
	// 		Id:             item.ID.String(),
	// 		Name:           item.MunicipalityName,
	// 		MunicipalityId: item.MunicipalityId.String(),
	// 	})
	// }
	return res, nil
}

func (m *BusinessServer) UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest) (*pb.Business, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.businessService.UpdateBusiness(ctx, req, md)
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
	return res, nil
}

func (m *BusinessServer) Feed(ctx context.Context, req *pb.FeedRequest) (*pb.FeedResponse, error) {
	var invalidProvinceId, invalidMunicipalityId, invalidLocation, invalidSearchMunicipalityType *epb.BadRequest_FieldViolation
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
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		if invalidLocation != nil {
			st, _ = st.WithDetails(
				invalidLocation,
			)
		}
		if invalidSearchMunicipalityType != nil {
			st, _ = st.WithDetails(
				invalidSearchMunicipalityType,
			)
		}
		if invalidProvinceId != nil {
			st, _ = st.WithDetails(
				invalidProvinceId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.Feed(ctx, req, md)
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
	return res, nil
}

func (m *BusinessServer) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest) (*pb.GetBusinessResponse, error) {
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
	res, err := m.businessService.GetBusiness(ctx, req, md)
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
	return res, nil
}

func (m *BusinessServer) GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest) (*pb.GetBusinessWithDistanceResponse, error) {
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
	res, err := m.businessService.GetBusinessWithDistance(ctx, req, md)
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
	return res, nil
}
