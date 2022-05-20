package app

import (
	"context"

	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (m *ObjectStorageServer) GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutRequest) (*pb.GetPresignedPutResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	res, err := m.objectStorageService.GetPresignedPutObject(&dto.GetPresignedPutObjectRequest{Metadata: &md, PhotoType: req.PhotoType.String(), LowQualityPhotoObject: req.LowQualityPhoto, HighQualityPhotoObject: req.HighQualityPhoto, ThumbnailQualityPhotoObject: req.ThumbnailQualityPhoto})
	if err != nil {
		switch err.Error() {
		case "authorizationtoken not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.GetPresignedPutResponse{LowQualityPhotoPresignedPutUrl: res.LowQualityPhotoPresignedPutUrl, HighQualityPhotoPresignedPutUrl: res.HighQualityPhotoPresignedPutUrl, ThumbnailPresignedPutUrl: res.ThumbnailQualityPhotoPresignedPutUrl}, nil
}
