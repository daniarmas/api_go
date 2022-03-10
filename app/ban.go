package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (m *BanServer) GetBannedDevice(ctx context.Context, req *gp.Empty) (*pb.GetBannedDeviceResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.banService.GetBannedDevice(&md)
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	if res != nil {
		return &pb.GetBannedDeviceResponse{BanExpirationTime: timestamppb.New(res.BanExpirationTime), CreateTime: timestamppb.New(res.CreateTime)}, nil
	} else {
		return &pb.GetBannedDeviceResponse{}, nil
	}
}

func (m *BanServer) GetBannedUser(ctx context.Context, req *gp.Empty) (*pb.GetBannedUserResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.banService.GetBannedUser(&md)
	if err != nil {
		switch err.Error() {
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	if res != nil {
		return &pb.GetBannedUserResponse{BanExpirationTime: timestamppb.New(res.BanExpirationTime), CreateTime: timestamppb.New(res.CreateTime)}, nil
	} else {
		return &pb.GetBannedUserResponse{}, nil
	}
}
