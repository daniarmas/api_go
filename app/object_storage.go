package app

import (
	"context"

	pb "github.com/daniarmas/api_go/pkg/grpc"
	utils "github.com/daniarmas/api_go/utils"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *ObjectStorageServer) GetPresignedPutObject(ctx context.Context, req *pb.GetPresignedPutObjectRequest) (*pb.GetPresignedPutObjectResponse, error) {
	var invalidThumbnailQualityPhoto, invalidHighQualityPhoto, invalidLowQualityPhoto *epb.BadRequest_FieldViolation
	var invalidArgs bool
	var st *status.Status
	md := utils.GetMetadata(ctx)
	if md.Authorization == nil {
		st = status.New(codes.Unauthenticated, "Unauthenticated user")
		return nil, st.Err()
	}
	if req.HighQualityPhoto == "" {
		invalidArgs = true
		invalidHighQualityPhoto = &epb.BadRequest_FieldViolation{
			Field:       "HighQualityPhoto",
			Description: "The HighQualityPhoto field is required",
		}
	}
	if req.LowQualityPhoto == "" {
		invalidArgs = true
		invalidLowQualityPhoto = &epb.BadRequest_FieldViolation{
			Field:       "LowQualityPhoto",
			Description: "The LowQualityPhoto field is required",
		}
	}
	if req.ThumbnailQualityPhoto == "" {
		invalidArgs = true
		invalidThumbnailQualityPhoto = &epb.BadRequest_FieldViolation{
			Field:       "ThumbnailQualityPhoto",
			Description: "The ThumbnailQualityPhoto field is required",
		}
	}
	if invalidArgs {
		st = status.New(codes.InvalidArgument, "Invalid Arguments")
		if invalidHighQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidHighQualityPhoto,
			)
		}
		if invalidLowQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidLowQualityPhoto,
			)
		}
		if invalidThumbnailQualityPhoto != nil {
			st, _ = st.WithDetails(
				invalidThumbnailQualityPhoto,
			)
		}
		return nil, st.Err()
	}
	res, err := m.objectStorageService.GetPresignedPutObject(ctx, req, md)
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
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return res, nil
}
