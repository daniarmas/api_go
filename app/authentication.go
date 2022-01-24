package app

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	ut "github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

// func naming(s string) string {
// 	if s == "foo" {
// 		return "Foo"
// 	}
// 	return s
// }

func (m *AuthenticationServer) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	verificationCode := models.VerificationCode{Code: ut.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceId: md.Get("deviceid")[0], CreateTime: time.Now(), UpdateTime: time.Now()}
	err := m.authenticationService.CreateVerificationCode(&verificationCode)
	if err != nil {
		switch err.Error() {
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
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	_, err := m.authenticationService.GetVerificationCode(&models.VerificationCode{Code: req.Code, Email: req.Email, Type: req.Type.String(), DeviceId: md.Get("deviceid")[0]}, &[]string{"id"})
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
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.SignIn(&models.VerificationCode{Code: req.Code, Email: req.Email, Type: "SignIn", DeviceId: md.Get("deviceid")[0]}, &md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.SignInResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken, User: &pb.User{Id: result.User.ID.String(), FullName: result.User.FullName, Alias: result.User.Alias, HighQualityPhoto: result.User.HighQualityPhoto, HighQualityPhotoBlurHash: result.User.HighQualityPhotoBlurHash, LowQualityPhoto: result.User.LowQualityPhoto, LowQualityPhotoBlurHash: result.User.LowQualityPhotoBlurHash, Thumbnail: result.User.Thumbnail, ThumbnailBlurHash: result.User.ThumbnailBlurHash, UserAddress: nil, Email: result.User.Email}}, nil
}

func (m *AuthenticationServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.SignUp(&req.FullName, &req.Alias, &models.VerificationCode{Code: req.Code, Email: req.Email, Type: "SignIn", DeviceId: md.Get("deviceid")[0]}, &md)
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
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.SignUpResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken, User: &pb.User{Id: result.User.ID.String(), FullName: result.User.FullName, Alias: result.User.Alias, HighQualityPhoto: result.User.HighQualityPhoto, HighQualityPhotoBlurHash: result.User.HighQualityPhotoBlurHash, LowQualityPhoto: result.User.LowQualityPhoto, LowQualityPhotoBlurHash: result.User.LowQualityPhotoBlurHash, Thumbnail: result.User.Thumbnail, ThumbnailBlurHash: result.User.ThumbnailBlurHash, UserAddress: nil, Email: result.User.Email}}, nil
}

func (m *AuthenticationServer) UserExists(ctx context.Context, req *pb.UserExistsRequest) (*gp.Empty, error) {
	var st *status.Status
	err := m.authenticationService.UserExists(&req.Alias)
	if err != nil {
		switch err.Error() {
		case "user already exists":
			st = status.New(codes.AlreadyExists, "User already exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
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
			st = status.New(codes.PermissionDenied, "User Banned")
		case "device banned":
			st = status.New(codes.PermissionDenied, "Device Banned")
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
	err := m.authenticationService.SignOut(&req.All, &req.AuthorizationTokenFk, &md)
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
