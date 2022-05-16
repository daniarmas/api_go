package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/utils"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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
			Address:        e.Address,
			Number:         e.Number,
			Coordinates:    &pb.Point{Latitude: e.Coordinates.Coords()[0], Longitude: e.Coordinates.Coords()[1]},
			Instructions:   e.Instructions,
			UserId:         e.UserId.String(),
			ProvinceId:     e.ProvinceId.String(),
			MunicipalityId: e.MunicipalityId.String(),
			CreateTime:     timestamppb.New(e.CreateTime),
			UpdateTime:     timestamppb.New(e.UpdateTime),
		})
	}
	return &pb.GetUserResponse{User: &pb.User{
		Id:                       getUserResponse.ID.String(),
		FullName:                 getUserResponse.FullName,
		HighQualityPhoto:         getUserResponse.HighQualityPhoto,
		HighQualityPhotoBlurHash: getUserResponse.HighQualityPhotoBlurHash,
		LowQualityPhoto:          getUserResponse.LowQualityPhoto,
		LowQualityPhotoBlurHash:  getUserResponse.LowQualityPhotoBlurHash,
		Thumbnail:                getUserResponse.Thumbnail,
		ThumbnailBlurHash:        getUserResponse.ThumbnailBlurHash,
		Email:                    getUserResponse.Email,
		UserAddress:              userAddress,
		CreateTime:               timestamppb.New(getUserResponse.CreateTime),
		UpdateTime:               timestamppb.New(getUserResponse.UpdateTime),
	}}, nil
}

func (m *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	updateUserResponse, err := m.userService.UpdateUser(&dto.UpdateUserRequest{Metadata: &md, Email: req.Email, FullName: req.FullName, Thumbnail: req.Thumbnail, ThumbnailBlurHash: req.ThumbnailBlurHash, HighQualityPhoto: req.HighQualityPhoto, HighQualityPhotoBlurHash: req.HighQualityPhotoBlurHash, LowQualityPhoto: req.LowQualityPhoto, LowQualityPhotoBlurHash: req.LowQualityPhotoBlurHash, Code: req.Code})
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
		case "user already exist":
			st = status.New(codes.AlreadyExists, "User already exists")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	// userAddress := make([]*pb.UserAddress, 0, len(updateUserResponse.UserAddress))
	// for _, e := range updateUserResponse.UserAddress {
	// 	userAddress = append(userAddress, &pb.UserAddress{
	// 		Id:             e.ID.String(),
	// 		Tag:            e.Tag,
	// 		ResidenceType:  *utils.ParseResidenceType(e.ResidenceType),
	// 		BuildingNumber: e.BuildingNumber,
	// 		HouseNumber:    e.HouseNumber,
	// 		Coordinates:    &pb.Point{Latitude: e.Coordinates.Coords()[0], Longitude: e.Coordinates.Coords()[1]},
	// 		Description:    e.Description,
	// 		UserId:         e.UserId.String(),
	// 		ProvinceId:     e.ProvinceId.String(),
	// 		MunicipalityId: e.MunicipalityId.String(),
	// 		CreateTime:     e.CreateTime.String(),
	// 		UpdateTime:     e.UpdateTime.String(),
	// 	})
	// }
	return &pb.UpdateUserResponse{User: &pb.User{
		Id:                       updateUserResponse.User.ID.String(),
		FullName:                 updateUserResponse.User.FullName,
		HighQualityPhoto:         updateUserResponse.User.HighQualityPhoto,
		HighQualityPhotoBlurHash: updateUserResponse.User.HighQualityPhotoBlurHash,
		LowQualityPhoto:          updateUserResponse.User.LowQualityPhoto,
		LowQualityPhotoBlurHash:  updateUserResponse.User.LowQualityPhotoBlurHash,
		Thumbnail:                updateUserResponse.User.Thumbnail,
		ThumbnailBlurHash:        updateUserResponse.User.ThumbnailBlurHash,
		Email:                    updateUserResponse.User.Email,
		// UserAddress:              userAddress,
		CreateTime: timestamppb.New(updateUserResponse.User.CreateTime),
		UpdateTime: timestamppb.New(updateUserResponse.User.UpdateTime),
	}}, nil
}
