package app

import (
	"context"
	"net/mail"
	"strconv"

	pb "github.com/daniarmas/api_go/pkg"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "record not found":
			st = status.New(codes.NotFound, "Not found")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
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
	res, err := m.authenticationService.SignIn(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
	return res, nil
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
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
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
	md := utils.GetMetadata(ctx)
	result, err := m.authenticationService.CheckSession(ctx, md)
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
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "user banned":
			st = status.New(codes.PermissionDenied, "User banned")
		case "device banned":
			st = status.New(codes.PermissionDenied, "Device banned")
		case "app banned":
			st = status.New(codes.PermissionDenied, "App banned")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token expired":
			st = status.New(codes.Unauthenticated, "Authorization token expired")
		case "authorization token contains an invalid number of segments", "authorization token signature is invalid":
			st = status.New(codes.Unauthenticated, "Authorization token invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.CheckSessionResponse{IpAddresses: *result}, nil
}

func (m *AuthenticationServer) SignOut(ctx context.Context, req *pb.SignOutRequest) (*gp.Empty, error) {
	var invalidAuthorizationTokenId *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.AuthorizationTokenId != "" {
		if !utils.IsValidUUID(&req.AuthorizationTokenId) {
			invalidArgs = true
			invalidAuthorizationTokenId = &epb.BadRequest_FieldViolation{
				Field:       "AuthorizationTokenId",
				Description: "The AuthorizationTokenId field is not a valid uuid v4",
			}
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidAuthorizationTokenId != nil {
			st, _ = st.WithDetails(
				invalidAuthorizationTokenId,
			)
		}
		return nil, st.Err()
	}
	res, err := m.authenticationService.SignOut(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *AuthenticationServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	var invalidRefreshToken *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.RefreshToken == "" {
		invalidArgs = true
		invalidRefreshToken = &epb.BadRequest_FieldViolation{
			Field:       "RefreshToken",
			Description: "The RefreshToken field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidRefreshToken != nil {
			st, _ = st.WithDetails(
				invalidRefreshToken,
			)
		}
		return nil, st.Err()
	}
	res, err := m.authenticationService.RefreshToken(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "refreshtoken expired":
			st = status.New(codes.Unauthenticated, "RefreshToken expired")
		case "authorization token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
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

func (m *AuthenticationServer) ListSession(ctx context.Context, req *gp.Empty) (*pb.ListSessionResponse, error) {
	var st *status.Status
	md := utils.GetMetadata(ctx)
	res, err := m.authenticationService.ListSession(ctx, md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
