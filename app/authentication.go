package app

import (
	"context"
	"net/mail"
	"strconv"

	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *AuthenticationServer) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest) (*gp.Empty, error) {
	var (
		invalidEmail *epb.BadRequest_FieldViolation
		invalidType  *epb.BadRequest_FieldViolation
	)
	var invalidArgs bool
	var st *status.Status
	meta := utils.GetMetadata(ctx)
	if req.Email == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "Email",
			Description: "The email field is required",
		}
	} else {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			invalidArgs = true
			invalidEmail = &epb.BadRequest_FieldViolation{
				Field:       "Email",
				Description: "The email field is invalid",
			}
		}
	}
	if req.Type == pb.VerificationCodeType_VerificationCodeTypeUnspecified {
		invalidArgs = true
		invalidType = &epb.BadRequest_FieldViolation{
			Field:       "Type",
			Description: "The type field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidEmail != nil {
			st, _ = st.WithDetails(
				invalidEmail,
			)
		}
		if invalidType != nil {
			st, _ = st.WithDetails(
				invalidType,
			)
		}
		return nil, st.Err()
	}
	res, err := m.authenticationService.CreateVerificationCode(ctx, req, meta)
	if err != nil {
		switch err.Error() {
		case "banned user":
			st = status.New(codes.PermissionDenied, "User banned")
		case "banned device":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "app banned":
			st = status.New(codes.PermissionDenied, "App banned")
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

func (m *AuthenticationServer) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest) (*gp.Empty, error) {
	var (
		invalidEmail *epb.BadRequest_FieldViolation
		invalidType  *epb.BadRequest_FieldViolation
		invalidCode  *epb.BadRequest_FieldViolation
	)
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Email == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "Email",
			Description: "The email field is required",
		}
	} else {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			invalidArgs = true
			invalidEmail = &epb.BadRequest_FieldViolation{
				Field:       "Email",
				Description: "The email field is invalid",
			}
		}
	}
	if req.Code == "" {
		invalidArgs = true
		invalidCode = &epb.BadRequest_FieldViolation{
			Field:       "Code",
			Description: "The code field is required",
		}
	} else {
		if _, err := strconv.Atoi(req.Code); err != nil {
			invalidArgs = true
			invalidCode = &epb.BadRequest_FieldViolation{
				Field:       "Code",
				Description: "The code field is invalid",
			}
		}
	}
	if req.Type == pb.VerificationCodeType_VerificationCodeTypeUnspecified {
		invalidArgs = true
		invalidType = &epb.BadRequest_FieldViolation{
			Field:       "Type",
			Description: "The type field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
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
		if invalidType != nil {
			st, _ = st.WithDetails(
				invalidType,
			)
		}
		return nil, st.Err()
	}
	res, err := m.authenticationService.GetVerificationCode(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "record not found":
			st = status.New(codes.NotFound, "Not found")
			return nil, st.Err()
		default:
			st = status.New(codes.Internal, "Internal server error")
			return nil, st.Err()
		}
	}
	return res, nil
}

func (m *AuthenticationServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	var (
		invalidEmail *epb.BadRequest_FieldViolation
		invalidCode  *epb.BadRequest_FieldViolation
	)
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Email == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "Email",
			Description: "The email field is required",
		}
	} else {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			invalidArgs = true
			invalidEmail = &epb.BadRequest_FieldViolation{
				Field:       "Email",
				Description: "The email field is invalid",
			}
		}
	}
	if req.Code == "" {
		invalidArgs = true
		invalidCode = &epb.BadRequest_FieldViolation{
			Field:       "Code",
			Description: "The code field is required",
		}
	} else {
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
	res, err := m.authenticationService.SignIn(&models.VerificationCode{Code: req.Code, Email: req.Email, Type: "SignIn", DeviceIdentifier: *md.DeviceIdentifier}, md)
	if err != nil {
		switch err.Error() {
		case "verification code not found":
			st = status.New(codes.NotFound, "VerificationCode Not found")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		case "user banned":
			st = status.New(codes.PermissionDenied, "User banned")
		case "device banned":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "app banned":
			st = status.New(codes.PermissionDenied, "App banned")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	permissions := make([]*pb.Permission, 0, len(res.User.UserPermissions))
	for _, item := range res.User.UserPermissions {
		permissions = append(permissions, &pb.Permission{
			Id:         item.ID.String(),
			Name:       item.Name,
			UserId:     item.UserId.String(),
			BusinessId: item.BusinessId.String(),
			CreateTime: timestamppb.New(item.CreateTime),
			UpdateTime: timestamppb.New(item.UpdateTime),
		})
	}
	return &pb.SignInResponse{RefreshToken: res.RefreshToken, AuthorizationToken: res.AuthorizationToken, User: &pb.User{Id: res.User.ID.String(), FullName: res.User.FullName, HighQualityPhoto: res.User.HighQualityPhoto, HighQualityPhotoBlurHash: res.User.HighQualityPhotoBlurHash, LowQualityPhoto: res.User.LowQualityPhoto, LowQualityPhotoBlurHash: res.User.LowQualityPhotoBlurHash, Thumbnail: res.User.Thumbnail, ThumbnailBlurHash: res.User.ThumbnailBlurHash, UserAddress: nil, Email: res.User.Email, Permissions: permissions, CreateTime: timestamppb.New(res.User.CreateTime), UpdateTime: timestamppb.New(res.User.UpdateTime)}}, nil
}

func (m *AuthenticationServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var (
		invalidEmail    *epb.BadRequest_FieldViolation
		invalidCode     *epb.BadRequest_FieldViolation
		invalidFullname *epb.BadRequest_FieldViolation
	)
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.Email == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "Email",
			Description: "The email field is required",
		}
	} else {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			invalidArgs = true
			invalidEmail = &epb.BadRequest_FieldViolation{
				Field:       "Email",
				Description: "The email field is invalid",
			}
		}
	}
	if req.FullName == "" {
		invalidArgs = true
		invalidFullname = &epb.BadRequest_FieldViolation{
			Field:       "Fullname",
			Description: "The fullname field is required",
		}
	}
	if req.Code == "" {
		invalidArgs = true
		invalidCode = &epb.BadRequest_FieldViolation{
			Field:       "Code",
			Description: "The code field is required",
		}
	} else {
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
		if invalidFullname != nil {
			st, _ = st.WithDetails(
				invalidFullname,
			)
		}
		return nil, st.Err()
	}
	res, err := m.authenticationService.SignUp(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "verification code not found":
			st = status.New(codes.NotFound, "VerificationCode Not found")
		case "user exists":
			st = status.New(codes.AlreadyExists, "User exists")
		case "user banned":
			st = status.New(codes.PermissionDenied, "User banned")
		case "device banned":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "app banned":
			st = status.New(codes.PermissionDenied, "App banned")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *AuthenticationServer) CheckSession(ctx context.Context, req *gp.Empty) (*pb.CheckSessionResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.CheckSession(&md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "user banned":
			st = status.New(codes.PermissionDenied, "User banned")
		case "device banned":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "app banned":
			st = status.New(codes.PermissionDenied, "App banned")
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
	return &pb.CheckSessionResponse{IpAddresses: *result}, nil
}

func (m *AuthenticationServer) SignOut(ctx context.Context, req *pb.SignOutRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	err := m.authenticationService.SignOut(&req.All, &req.AuthorizationTokenId, &md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
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
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.RefreshToken(&req.RefreshToken, &md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "refreshtoken expired":
			st = status.New(codes.Unauthenticated, "RefreshToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.RefreshTokenResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken}, nil
}

func (m *AuthenticationServer) ListSession(ctx context.Context, req *pb.ListSessionRequest) (*pb.ListSessionResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.ListSession(&md)
	if err != nil {
		switch err.Error() {
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	sessions := make([]*pb.Session, 0, len(*result.Sessions))
	for _, e := range *result.Sessions {
		var actual bool = false
		if e.Device.ID == result.ActualDeviceId {
			actual = true
		}
		sessions = append(sessions, &pb.Session{
			Id:            e.ID.String(),
			Platform:      *utils.ParsePlatformType(&e.Platform),
			SystemVersion: e.SystemVersion,
			Model:         e.Model,
			App:           *utils.ParseAppType(&e.App),
			AppVersion:    e.AppVersion,
			DeviceId:      e.DeviceId.String(),
			Actual:        actual,
		})
	}
	return &pb.ListSessionResponse{Sessions: sessions}, nil
}
