package app

import (
	"context"
	"net/mail"
	"strconv"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *UserServer) GetAddressInfo(ctx context.Context, req *pb.GetAddressInfoRequest) (*pb.GetAddressInfoResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	getAddressInfoRes, err := m.userService.GetAddressInfo(&dto.GetAddressInfoRequest{Metadata: &md, Coordinates: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}})
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
	return &pb.GetAddressInfoResponse{ProvinceId: getAddressInfoRes.ProvinceId.String(), ProvinceName: getAddressInfoRes.ProvinceName, MunicipalityName: getAddressInfoRes.MunicipalityName, MunicipalityId: getAddressInfoRes.MunicipalityId.String(), ProvinceNameAbbreviation: getAddressInfoRes.ProvinceNameAbbreviation}, nil
}

func (m *UserServer) GetUser(ctx context.Context, req *gp.Empty) (*pb.GetUserResponse, error) {
	var st *status.Status
	meta := utils.GetMetadata(ctx)
	res, err := m.userService.GetUser(ctx, meta)
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
