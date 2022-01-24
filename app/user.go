package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *UserServer) GetUser(ctx context.Context, req *gp.Empty) (*pb.GetUserResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	user, err := m.userService.GetUser(&md)
	if err != nil {
		switch err.Error() {
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
	return &pb.GetUserResponse{User: &pb.User{
		Id:                       user.ID.String(),
		FullName:                 user.FullName,
		Alias:                    user.Alias,
		HighQualityPhoto:         user.HighQualityPhoto,
		HighQualityPhotoBlurHash: user.HighQualityPhotoBlurHash,
		LowQualityPhoto:          user.LowQualityPhoto,
		LowQualityPhotoBlurHash:  user.LowQualityPhotoBlurHash,
		Thumbnail:                user.Thumbnail,
		ThumbnailBlurHash:        user.ThumbnailBlurHash,
		Email:                    user.Email,
		CreateTime:               user.CreateTime.String(),
		UpdateTime:               user.UpdateTime.String(),
	}}, nil
}
