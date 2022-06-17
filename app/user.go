package app

import (
	"context"
	"net/mail"
	"strconv"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *UserServer) GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest) (*pb.GetAddressInfoResponse, error) {
	var invalidProvinceId, invalidMunicipalityId, invalidLocation, invalidSearchMunicipalityType *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
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
	res, err := m.userService.GetAddressInfo(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "municipality not found":
			st = status.New(codes.InvalidArgument, "Location not available")
		case "province not found":
			st = status.New(codes.InvalidArgument, "Location not available")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *UserServer) GetUser(ctx context.Context, req *gp.Empty) (*pb.GetUserResponse, error) {
	var st *status.Status
	meta := utils.GetMetadata(ctx)
	if meta.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.userService.GetUser(ctx, meta)
	if err != nil {
		switch err.Error() {
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *UserServer) ListUserAddress(ctx context.Context, req *gp.Empty) (*pb.ListUserAddressResponse, error) {
	var st *status.Status
	meta := utils.GetMetadata(ctx)
	if meta.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	res, err := m.userService.ListUserAddress(ctx, req, meta)
	if err != nil {
		switch err.Error() {
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var invalidCode, invalidEmail, invalidId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated")
		return nil, st.Err()
	}
	if req.User.Email != "" {
		_, err := mail.ParseAddress(req.User.Email)
		if err != nil {
			invalidArgs = true
			invalidEmail = &epb.BadRequest_FieldViolation{
				Field:       "User.Email",
				Description: "The User.Email field is invalid",
			}
		}
	}
	if req.User.Id != "" {
		if !utils.IsValidUUID(&req.User.Id) {
			invalidArgs = true
			invalidId = &epb.BadRequest_FieldViolation{
				Field:       "User.Id",
				Description: "The User.Id field is not a valid uuid v4",
			}
		}
	}
	if req.Code != "" {
		if _, err := strconv.Atoi(req.Code); err != nil {
			invalidArgs = true
			invalidCode = &epb.BadRequest_FieldViolation{
				Field:       "Code",
				Description: "The code field is invalid",
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
		if invalidEmail != nil {
			st, _ = st.WithDetails(
				invalidEmail,
			)
		}
		if invalidCode != nil {
			st, _ = st.WithDetails(
				invalidCode,
			)
		}
		return nil, st.Err()
	}
	res, err := m.userService.UpdateUser(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "missing code":
			invalidCode := &epb.BadRequest_FieldViolation{
				Field:       "Code",
				Description: "The code field is required for update the email",
			}
			st = status.New(codes.InvalidArgument, "Invalid Arguments")
			st, _ = st.WithDetails(
				invalidCode,
			)
		case "not have permission":
			st = status.New(codes.PermissionDenied, "PermissionDenied")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "user already exist":
			st = status.New(codes.AlreadyExists, "User already exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
