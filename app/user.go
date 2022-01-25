package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *UserServer) GetUser(ctx context.Context, req *gp.Empty) (*pb.GetUserResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	getUserResponse, err := m.userService.GetUser(&md)
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
	userAddress := make([]*pb.UserAddress, 0, len(getUserResponse.UserAddress))
	for _, e := range getUserResponse.UserAddress {
		userAddress = append(userAddress, &pb.UserAddress{
			Id:             e.ID.String(),
			Tag:            e.Tag,
			ResidenceType:  *utils.ParseResidenceType(e.ResidenceType),
			BuildingNumber: e.BuildingNumber,
			HouseNumber:    e.HouseNumber,
			Coordinates:    &pb.Point{Latitude: e.Coordinates.Coords()[0], Longitude: e.Coordinates.Coords()[1]},
			Description:    e.Description,
			UserFk:         e.UserFk.String(),
			ProvinceFk:     e.ProvinceFk.String(),
			MunicipalityFk: e.MunicipalityFk.String(),
			CreateTime:     e.CreateTime.String(),
			UpdateTime:     e.UpdateTime.String(),
		})
	}
	return &pb.GetUserResponse{User: &pb.User{
		Id:                       getUserResponse.ID.String(),
		FullName:                 getUserResponse.FullName,
		Alias:                    getUserResponse.Alias,
		HighQualityPhoto:         getUserResponse.HighQualityPhoto,
		HighQualityPhotoBlurHash: getUserResponse.HighQualityPhotoBlurHash,
		LowQualityPhoto:          getUserResponse.LowQualityPhoto,
		LowQualityPhotoBlurHash:  getUserResponse.LowQualityPhotoBlurHash,
		Thumbnail:                getUserResponse.Thumbnail,
		ThumbnailBlurHash:        getUserResponse.ThumbnailBlurHash,
		Email:                    getUserResponse.Email,
		UserAddress:              userAddress,
		CreateTime:               getUserResponse.CreateTime.String(),
		UpdateTime:               getUserResponse.UpdateTime.String(),
	}}, nil
}
