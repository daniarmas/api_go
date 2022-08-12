package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	utils "github.com/daniarmas/api_go/utils"
	// epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *PermissionServer) ListPermission(ctx context.Context, req *pb.ListPermissionRequest) (*pb.ListPermissionResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
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
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

// func (m *PermissionServer) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*gp.Empty, error) {
// 	var invalidId *epb.BadRequest_FieldViolation
// 	var invalidArgs bool
// 	var st *status.Status
// 	md := utils.GetMetadata(ctx)
// 	if md.Authorization == nil {
// 		st = status.New(codes.Unauthenticated, "Unauthenticated")
// 		return nil, st.Err()
// 	}
// 	if req.Id == "" {
// 		invalidArgs = true
// 		invalidId = &epb.BadRequest_FieldViolation{
// 			Field:       "id",
// 			Description: "The id field is required",
// 		}
// 	} else if req.Id != "" {
// 		if !utils.IsValidUUID(&req.Id) {
// 			invalidArgs = true
// 			invalidId = &epb.BadRequest_FieldViolation{
// 				Field:       "id",
// 				Description: "The id field is not a valid uuid v4",
// 			}
// 		}
// 	}
// 	if invalidArgs {
// 		st = status.New(codes.InvalidArgument, "Invalid Arguments")
// 		if invalidId != nil {
// 			st, _ = st.WithDetails(
// 				invalidId,
// 			)
// 		}
// 		return nil, st.Err()
// 	}
// 	res, err := m.PermissionsService.DeletePermissions(ctx, req, md)
// 	if err != nil {
// 		switch err.Error() {
// 		case "unauthenticated Permissions":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated Permissions")
// 		case "access token contains an invalid number of segments", "access token signature is invalid":
// 			st = status.New(codes.Unauthenticated, "Access token is invalid")
// 		case "access token expired":
// 			st = status.New(codes.Unauthenticated, "Access token is expired")
// 		case "unauthenticated":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated")
// 		case "Permissions not found":
// 			st = status.New(codes.NotFound, "Permissions not found")
// 		case "authorization token not found":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated")
// 		case "authorization token expired":
// 			st = status.New(codes.Unauthenticated, "Authorization token expired")
// 		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
// 			st = status.New(codes.Unauthenticated, "Authorization token invalid")
// 		case "permission denied":
// 			st = status.New(codes.PermissionDenied, "Permission denied")
// 		default:
// 			st = status.New(codes.Internal, "Internal server error")
// 		}
// 		return nil, st.Err()
// 	}
// 	return res, nil
// }

// func (m *PermissionsServer) CreatePermissions(ctx context.Context, req *pb.CreatePermissionsRequest) (*pb.Permissions, error) {
// 	var invalidName, invalidVersion *epb.BadRequest_FieldViolation
// 	var invalidArgs bool
// 	var st *status.Status
// 	md := utils.GetMetadata(ctx)
// 	if req.Permissions.Name == "" {
// 		invalidArgs = true
// 		invalidName = &epb.BadRequest_FieldViolation{
// 			Field:       "Permissions.Name",
// 			Description: "The Permissions.Name field is required",
// 		}
// 	}
// 	if req.Permissions.Version == "" {
// 		invalidArgs = true
// 		invalidVersion = &epb.BadRequest_FieldViolation{
// 			Field:       "Permissions.Version",
// 			Description: "The Permissions.Version field is required",
// 		}
// 	} else if req.Permissions.Version != "" {
// 		if !utils.RegexpSemanticVersion(&req.Permissions.Version) {
// 			invalidArgs = true
// 			invalidVersion = &epb.BadRequest_FieldViolation{
// 				Field:       "Permissions.Version",
// 				Description: "The Permissions.Version field value is not a valid semantic version",
// 			}
// 		}
// 	}
// 	if invalidArgs {
// 		st = status.New(codes.InvalidArgument, "Invalid Arguments")
// 		if invalidName != nil {
// 			st, _ = st.WithDetails(
// 				invalidName,
// 			)
// 		}
// 		if invalidVersion != nil {
// 			st, _ = st.WithDetails(
// 				invalidVersion,
// 			)
// 		}
// 		return nil, st.Err()
// 	}
// 	res, err := m.PermissionsService.CreatePermissions(ctx, req, md)
// 	if err != nil {
// 		switch err.Error() {
// 		case "unauthenticated Permissions":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated Permissions")
// 		case "access token contains an invalid number of segments", "access token signature is invalid":
// 			st = status.New(codes.Unauthenticated, "Access token is invalid")
// 		case "access token expired":
// 			st = status.New(codes.Unauthenticated, "Access token is expired")
// 		case "unauthenticated":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated")
// 		case "authorization token not found":
// 			st = status.New(codes.Unauthenticated, "Unauthenticated")
// 		case "authorization token expired":
// 			st = status.New(codes.Unauthenticated, "Authorization token expired")
// 		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
// 			st = status.New(codes.Unauthenticated, "Authorization token invalid")
// 		case "permission denied":
// 			st = status.New(codes.PermissionDenied, "Permission denied")
// 		default:
// 			st = status.New(codes.Internal, "Internal server error")
// 		}
// 		return nil, st.Err()
// 	}
// 	return res, nil
// }
