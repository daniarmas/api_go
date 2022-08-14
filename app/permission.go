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

func (m *PermissionServer) ListPermission(ctx context.Context, req *pb.ListPermissionRequest) (*pb.ListPermissionResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	res, err := m.permissionService.ListPermission(ctx, req, md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *PermissionServer) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*gp.Empty, error) {
	var invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
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
	res, err := m.permissionService.DeletePermission(ctx, req, md)
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
		case "permission not found":
			st = status.New(codes.NotFound, "Permission not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *PermissionServer) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.Permission, error) {
	var invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
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
	res, err := m.permissionService.GetPermission(ctx, req, md)
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
		case "permission not found":
			st = status.New(codes.NotFound, "Permission not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *PermissionServer) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.Permission, error) {
	var invalidName *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.Permission.Name == "" {
		invalidArgs = true
		invalidName = &epb.BadRequest_FieldViolation{
			Field:       "permission.Name",
			Description: "The permission.Name field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		return nil, st.Err()
	}
	res, err := m.permissionService.CreatePermission(ctx, req, md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
