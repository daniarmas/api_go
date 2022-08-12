package app

import (
	"context"
	"net/mail"
	"strconv"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *AuthenticationServer) SessionExists(ctx context.Context, req *pb.SessionExistsRequest) (*pb.SessionExistsResponse, error) {
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
	res, err := m.authenticationService.SessionExists(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "verification code not found":
			st = status.New(codes.NotFound, "Verification code not found")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		case "session not exists":
			st = status.New(codes.NotFound, "Session not exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

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
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
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
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "verification code not found":
			st = status.New(codes.NotFound, "Verification code not found")
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
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "verification code not found":
			st = status.New(codes.NotFound, "Verification code not found")
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		case "session limit reached":
			st = status.New(codes.PermissionDenied, "Session limit reached")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}

func (m *AuthenticationServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var (
		invalidEmail          *epb.BadRequest_FieldViolation
		invalidCode           *epb.BadRequest_FieldViolation
		invalidFullname       *epb.BadRequest_FieldViolation
		invalidName           *epb.BadRequest_FieldViolation
		invalidAddress        *epb.BadRequest_FieldViolation
		invalidNumber         *epb.BadRequest_FieldViolation
		invalidCoordinates    *epb.BadRequest_FieldViolation
		invalidProvinceId     *epb.BadRequest_FieldViolation
		invalidMunicipalityId *epb.BadRequest_FieldViolation
	)
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if req.UserAddress.Coordinates == nil {
		invalidArgs = true
		invalidCoordinates = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.Coordinates",
			Description: "The userAddress.Coordinates field is required",
		}
	} else if req.UserAddress.Coordinates != nil {
		if req.UserAddress.Coordinates.Latitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "userAddress.Coordinates.Latitude",
				Description: "The userAddress.Coordinates.Latitude field is required",
			}
		} else if req.UserAddress.Coordinates.Longitude == 0 {
			invalidArgs = true
			invalidCoordinates = &epb.BadRequest_FieldViolation{
				Field:       "userAddress.Coordinates.Longitude",
				Description: "The userAddress.Coordinates.Longitude field is required",
			}
		}
	}
	if req.UserAddress.ProvinceId == "" {
		invalidArgs = true
		invalidProvinceId = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.provinceId",
			Description: "The userAddress.provinceId field is required",
		}
	} else if req.UserAddress.ProvinceId != "" {
		if !utils.IsValidUUID(&req.UserAddress.ProvinceId) {
			invalidArgs = true
			invalidProvinceId = &epb.BadRequest_FieldViolation{
				Field:       "userAddress.provinceId",
				Description: "The userAddress.provinceId field is not a valid uuid v4",
			}
		}
	}
	if req.UserAddress.MunicipalityId == "" {
		invalidArgs = true
		invalidMunicipalityId = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.municipalityId",
			Description: "The userAddress.municipalityId field is required",
		}
	} else if req.UserAddress.MunicipalityId != "" {
		if !utils.IsValidUUID(&req.UserAddress.MunicipalityId) {
			invalidArgs = true
			invalidMunicipalityId = &epb.BadRequest_FieldViolation{
				Field:       "userAddress.municipalityId",
				Description: "The userAddress.municipalityId field is not a valid uuid v4",
			}
		}
	}
	if req.UserAddress.Name == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.name",
			Description: "The userAddress.name field is required",
		}
	}
	if req.UserAddress.Address == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.address",
			Description: "The userAddress.address field is required",
		}
	}
	if req.UserAddress.Number == "" {
		invalidArgs = true
		invalidEmail = &epb.BadRequest_FieldViolation{
			Field:       "userAddress.number",
			Description: "The userAddress.number field is required",
		}
	}
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
		if invalidAddress != nil {
			st, _ = st.WithDetails(
				invalidAddress,
			)
		}
		if invalidCoordinates != nil {
			st, _ = st.WithDetails(
				invalidCoordinates,
			)
		}
		if invalidName != nil {
			st, _ = st.WithDetails(
				invalidName,
			)
		}
		if invalidNumber != nil {
			st, _ = st.WithDetails(
				invalidNumber,
			)
		}
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
		return nil, st.Err()
	}
	res, err := m.authenticationService.SignUp(ctx, req, md)
	if err != nil {
		switch err.Error() {
		case "access token contains an invalid number of segments", "access token signature is invalid":
			st = status.New(codes.Unauthenticated, "Access token is invalid")
		case "access token expired":
			st = status.New(codes.Unauthenticated, "Access token is expired")
		case "unauthenticated application":
			st = status.New(codes.Unauthenticated, "Unauthenticated application")
		case "verification code not found":
			st = status.New(codes.NotFound, "Verification code not found")
		case "user exists":
			st = status.New(codes.AlreadyExists, "User already exists")
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
	res, err := m.authenticationService.CheckSession(ctx, md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
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
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
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
		case "refresh token not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "refresh token expired":
			st = status.New(codes.Unauthenticated, "Refresh token expired")
		case "refresh token contains an invalid number of segments", "refresh token signature is invalid":
			st = status.New(codes.Unauthenticated, "Refresh token invalid")
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
		case "unauthenticated user":
			st = status.New(codes.Unauthenticated, "Unauthenticated user")
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
