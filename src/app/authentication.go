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

func (m *AuthenticationServer) CreateVerificationCode(ctx context.Context, req *pb.CreateVerificationCodeRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	verificationCode := datastruct.VerificationCode{Code: ut.EncodeToString(6), Email: req.Email, Type: req.Type.Enum().String(), DeviceId: md.Get("deviceid")[0], CreateTime: time.Now(), UpdateTime: time.Now()}
	err := m.authenticationService.CreateVerificationCode(&verificationCode)
	if err != nil {
		switch err.Error() {
		case "banned user":
			st = status.New(codes.PermissionDenied, "Banned user")
		case "banned device":
			st = status.New(codes.PermissionDenied, "Banned device ")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) GetVerificationCode(ctx context.Context, req *pb.GetVerificationCodeRequest) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	err := m.authenticationService.GetVerificationCode(req.Code, req.Email, req.Type.String(), md.Get("deviceid")[0])
	if err != nil {
		switch err.Error() {
		case "user banned":
			st = status.New(codes.NotFound, "User banned")
		case "device banned":
			st = status.New(codes.NotFound, "Device banned")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}
