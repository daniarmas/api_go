package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *BusinessServer) ModifyBusinessRolePermission(ctx context.Context, req *pb.ModifyBusinessRolePermissionRequest) (*gp.Empty, error) {
	var invalidBusinessRoleId, invalidPermissionId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.BusinessRoleId == "" {
		invalidArgs = true
		invalidBusinessRoleId = &epb.BadRequest_FieldViolation{
			Field:       "id",
			Description: "The id field is required",
		}
	} else if req.BusinessRoleId != "" {
		if !utils.IsValidUUID(&req.BusinessRoleId) {
			invalidArgs = true
			invalidBusinessRoleId = &epb.BadRequest_FieldViolation{
				Field:       "businessRoleId",
				Description: "The businessRoleId field is not a valid uuid v4",
			}
		}
	}
	for _, i := range req.PermissionIds {
		if !utils.IsValidUUID(&i) {
			invalidArgs = true
			invalidPermissionId = &epb.BadRequest_FieldViolation{
				Field:       "permissionIds",
				Description: "The permissionIds field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidBusinessRoleId != nil {
			st, _ = st.WithDetails(
				invalidBusinessRoleId,
			)
		}
		if invalidPermissionId != nil {
			st, _ = st.WithDetails(
				invalidPermissionId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.ModifyBusinessRolePermission(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "business role not found":
			st = status.New(codes.AlreadyExists, "Business role not found")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) UpdateBusinessRole(ctx context.Context, req *pb.UpdateBusinessRoleRequest) (*pb.BusinessRole, error) {
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
	res, err := m.businessService.UpdateBusinessRole(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "business role not found":
			st = status.New(codes.AlreadyExists, "Business role not found")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) DeleteBusinessRole(ctx context.Context, req *pb.DeleteBusinessRoleRequest) (*gp.Empty, error) {
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
			Field:       "businessRole.businessId",
			Description: "The businessRole.businessId field is required",
		}
	} else if req.Id != "" {
		if !utils.IsValidUUID(&req.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "businessRole.businessId",
				Description: "The businessRole.businessId field is not a valid uuid v4",
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
	res, err := m.businessService.DeleteBusinessRole(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "already exists a business with that name":
			st = status.New(codes.AlreadyExists, "Already exists a business with that name")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) CreateBusinessRole(ctx context.Context, req *pb.CreateBusinessRoleRequest) (*pb.BusinessRole, error) {
	var invalidBusinessId, invalidName, invalidBusinessRolePermissionIds *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	for _, i := range req.BusinessRole.Permissions {
		if i.Id == "" {
			invalidArgs = true
			invalidBusinessRolePermissionIds = &epb.BadRequest_FieldViolation{
				Field:       "businessRole.permissions.id",
				Description: "The businessRole.permissions.id field is required",
			}
		} else if i.Id != "" {
			if !utils.IsValidUUID(&i.Id) {
				invalidArgs = true
				invalidBusinessRolePermissionIds = &epb.BadRequest_FieldViolation{
					Field:       "businessRole.permissions.id",
					Description: "The businessRole.permissions.id field is not a valid uuid v4",
				}
			}
		}
	}
	if req.BusinessRole.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "businessRole.name",
			Description: "The businessRole.name field is required",
		}
	}
	if req.BusinessRole.BusinessId == "" {
		invalidArgs = true
		invalidBusinessId = &epb.BadRequest_FieldViolation{
			Field:       "businessRole.businessId",
			Description: "The businessRole.businessId field is required",
		}
	} else if req.BusinessRole.BusinessId != "" {
		if !utils.IsValidUUID(&req.BusinessRole.BusinessId) {
			invalidArgs = true
			invalidBusinessId = &epb.BadRequest_FieldViolation{
				Field:       "businessRole.businessId",
				Description: "The businessRole.businessId field is not a valid uuid v4",
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
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		if invalidBusinessRolePermissionIds != nil {
			st, _ = st.WithDetails(
				invalidBusinessRolePermissionIds,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.CreateBusinessRole(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "already exists a business with that name":
			st = status.New(codes.AlreadyExists, "Already exists a business with that name")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) ListBusinessRole(ctx context.Context, req *pb.ListBusinessRoleRequest) (*pb.ListBusinessRoleResponse, error) {
	var invalidBusinessId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
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
		if invalidBusinessId != nil {
			st, _ = st.WithDetails(
				invalidBusinessId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.ListBusinessRole(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "already exists a business with that name":
			st = status.New(codes.AlreadyExists, "Already exists a business with that name")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) UpdatePartnerApplication(ctx context.Context, req *pb.UpdatePartnerApplicationRequest) (*pb.PartnerApplication, error) {
	var invalidBusinessName, invalidId, invalidDescription, invalidCoordinates, invalidMunicipalityId, invalidProvinceId *epb.BadRequest_FieldViolation
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
	if req.PartnerApplication.Coordinates == nil {
		invalidArgs = true
		invalidCoordinates = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.Coordinates",
			Description: "The PartnerApplication.Coordinates field is required",
		}
	} else if req.PartnerApplication.Coordinates != nil {
		if req.PartnerApplication.Coordinates.Latitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.Coordinates.Latitude",
				Description: "The PartnerApplication.Coordinates.Latitude field is required",
			}
		} else if req.PartnerApplication.Coordinates.Longitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.Coordinates.Longitude",
				Description: "The PartnerApplication.Coordinates.Longitude field is required",
			}
		}
	}
	if req.PartnerApplication.BusinessName == "" {
		invalidArgs = true
		invalidBusinessName = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.BusinessName",
			Description: "The PartnerApplication.BusinessName field is required",
		}
	}
	if req.PartnerApplication.Description == "" {
		invalidArgs = true
		invalidDescription = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.Description",
			Description: "The PartnerApplication.Description field is required",
		}
	}
	if req.PartnerApplication.ProvinceId == "" {
		invalidArgs = true
		invalidProvinceId = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.ProvinceId",
			Description: "The PartnerApplication.ProvinceId field is required",
		}
	} else if req.PartnerApplication.ProvinceId != "" {
		if !utils.IsValidUUID(&req.PartnerApplication.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.ProvinceId",
				Description: "The PartnerApplication.ProvinceId field is not a valid uuid v4",
			}
		}
	}
	if req.PartnerApplication.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.MunicipalityId",
			Description: "The PartnerApplication.MunicipalityId field is required",
		}
	} else if req.PartnerApplication.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.PartnerApplication.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.MunicipalityId",
				Description: "The PartnerApplication.MunicipalityId field is not a valid uuid v4",
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
		if invalidDescription != nil {
			st, _ = st.WithDetails(
				invalidDescription,
			)
		}
		if invalidId != nil {
			st, _ = st.WithDetails(
				invalidId,
			)
		}
		if invalidBusinessName != nil {
			st, _ = st.WithDetails(
				invalidBusinessName,
			)
		}
		if invalidCoordinates != nil {
			st, _ = st.WithDetails(
				invalidCoordinates,
			)
		}
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.UpdatePartnerApplication(ctx, req, md)
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
		case "partner application not found":
			st = status.New(codes.NotFound, "Partner application not found")
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "permission denied":
			st = status.New(codes.AlreadyExists, "Permission denied")
		case "already exists a business with that name":
			st = status.New(codes.AlreadyExists, "Already exists a business with that name")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) ListPartnerApplication(ctx context.Context, req *pb.ListPartnerApplicationRequest) (*pb.ListPartnerApplicationResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.businessService.ListPartnerApplication(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "not permission":
			st = status.New(codes.PermissionDenied, "Permission Denied")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) CreatePartnerApplication(ctx context.Context, req *pb.CreatePartnerApplicationRequest) (*pb.PartnerApplication, error) {
	var invalidBusinessName, invalidDescription, invalidCoordinates, invalidMunicipalityId, invalidProvinceId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.PartnerApplication.Coordinates == nil {
		invalidArgs = true
		invalidCoordinates = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.Coordinates",
			Description: "The PartnerApplication.Coordinates field is required",
		}
	} else if req.PartnerApplication.Coordinates != nil {
		if req.PartnerApplication.Coordinates.Latitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.Coordinates.Latitude",
				Description: "The PartnerApplication.Coordinates.Latitude field is required",
			}
		} else if req.PartnerApplication.Coordinates.Longitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.Coordinates.Longitude",
				Description: "The PartnerApplication.Coordinates.Longitude field is required",
			}
		}
	}
	if req.PartnerApplication.BusinessName == "" {
		invalidArgs = true
		invalidBusinessName = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.BusinessName",
			Description: "The PartnerApplication.BusinessName field is required",
		}
	}
	if req.PartnerApplication.Description == "" {
		invalidArgs = true
		invalidDescription = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.Description",
			Description: "The PartnerApplication.Description field is required",
		}
	}
	if req.PartnerApplication.ProvinceId == "" {
		invalidArgs = true
		invalidProvinceId = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.ProvinceId",
			Description: "The PartnerApplication.ProvinceId field is required",
		}
	} else if req.PartnerApplication.ProvinceId != "" {
		if !utils.IsValidUUID(&req.PartnerApplication.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.ProvinceId",
				Description: "The PartnerApplication.ProvinceId field is not a valid uuid v4",
			}
		}
	}
	if req.PartnerApplication.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "PartnerApplication.MunicipalityId",
			Description: "The PartnerApplication.MunicipalityId field is required",
		}
	} else if req.PartnerApplication.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.PartnerApplication.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "PartnerApplication.MunicipalityId",
				Description: "The PartnerApplication.MunicipalityId field is not a valid uuid v4",
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
		if invalidDescription != nil {
			st, _ = st.WithDetails(
				invalidDescription,
			)
		}
		if invalidBusinessName != nil {
			st, _ = st.WithDetails(
				invalidBusinessName,
			)
		}
		if invalidCoordinates != nil {
			st, _ = st.WithDetails(
				invalidCoordinates,
			)
		}
		if invalidMunicipalityId != nil {
			st, _ = st.WithDetails(
				invalidMunicipalityId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.businessService.CreatePartnerApplication(ctx, req, md)
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
		case "already register as business user":
			st = status.New(codes.AlreadyExists, "Already register as business user")
		case "already exists a business with that name":
			st = status.New(codes.AlreadyExists, "Already exists a business with that name")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest) (*pb.CreateBusinessResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.businessService.CreateBusiness(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "HighQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "HighQualityPhotoObject missing")
		case "LowQualityPhotoObject missing":
			st = status.New(codes.InvalidArgument, "LowQualityPhotoObject missing")
		case "ThumbnailObject missing":
			st = status.New(codes.InvalidArgument, "ThumbnailObject missing")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *BusinessServer) UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest) (*pb.Business, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.businessService.UpdateBusiness(ctx, req, md)
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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

func (m *BusinessServer) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest) (*pb.Business, error) {
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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

func (m *BusinessServer) GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest) (*pb.Business, error) {
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
