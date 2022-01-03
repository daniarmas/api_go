package app

import (
	"context"
	"time"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/src/datastruct"
	ut "github.com/daniarmas/api_go/src/utils"
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
	verificationCode := datastruct.VerificationCode{Code: ut.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceId: md.Get("deviceid")[0], CreateTime: time.Now(), UpdateTime: time.Now()}
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
	// FieldMask
	// userDst := &datastruct.VerificationCode{} // a struct to copy to
	// mask, _ := fieldmask_utils.MaskFromPaths(req.FieldMask.Paths, naming)
	// fields := strings.Split(mask.String(), ",")
	// fieldmask_utils.StructToStruct(mask, req.Email, userDst)
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.GetVerificationCode(&datastruct.VerificationCode{Code: req.Code, Email: req.Email, Type: req.Type.String(), DeviceId: md.Get("deviceid")[0]}, &[]string{"id"})
	if err != nil {
		st = status.New(codes.Internal, "Internal server error")
		return nil, st.Err()
	} else if result == nil {
		st = status.New(codes.NotFound, "Not found")
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	result, err := m.authenticationService.SignIn(&datastruct.VerificationCode{Code: req.Code, Email: req.Email, Type: "SignIn", DeviceId: md.Get("deviceid")[0]}, &md)
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
