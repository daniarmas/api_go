package utils

import (
	"context"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var st = status.New(codes.Unauthenticated, "Incorrect metadata")
	var invalidDeviceId *epb.ErrorInfo
	var invalidAccessToken *epb.ErrorInfo
	var invalidPlatform *epb.ErrorInfo
	var invalidSystemVersion *epb.ErrorInfo
	var invalidModel *epb.ErrorInfo
	var invalidArgs bool
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md.Get("device-id")) == 0 {
		invalidArgs = true
		invalidDeviceId = &epb.ErrorInfo{
			Reason: "device-id metadata missing",
		}
	}
	if len(md.Get("access-token")) == 0 {
		invalidArgs = true
		invalidAccessToken = &epb.ErrorInfo{
			Reason: "access-token metadata missing",
		}
	} else if md.Get("access-token")[0] != "O8pzXjp4QMk4cAD60dHeoOnxdVsDc9" {
		invalidArgs = true
		invalidAccessToken = &epb.ErrorInfo{
			Reason: "access-token is incorrect",
		}
	}
	if len(md.Get("platform")) == 0 {
		invalidArgs = true
		invalidPlatform = &epb.ErrorInfo{
			Reason: "platform metadata missing",
		}
	}
	if len(md.Get("system-version")) == 0 {
		invalidArgs = true
		invalidSystemVersion = &epb.ErrorInfo{
			Reason: "system-version metadata missing",
		}
	}
	if len(md.Get("model")) == 0 {
		invalidArgs = true
		invalidModel = &epb.ErrorInfo{
			Reason: "model metadata missing",
		}
	}
	if invalidArgs {
		st, _ = st.WithDetails(
			invalidDeviceId,
		)
		st, _ = st.WithDetails(
			invalidAccessToken,
		)
		st, _ = st.WithDetails(
			invalidPlatform,
		)
		st, _ = st.WithDetails(
			invalidSystemVersion,
		)
		st, _ = st.WithDetails(
			invalidModel,
		)
		return nil, st.Err()
	}
	m, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, err
}

func StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ss.Context())
	if len(md.Get("deviceid")) == 0 {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return st.Err()
	}
	if len(md.Get("accesstoken")) == 0 {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return st.Err()
	}
	if len(md.Get("platform")) == 0 {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return st.Err()
	}
	if len(md.Get("systemversion")) == 0 {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return st.Err()
	}
	if len(md.Get("model")) == 0 {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return st.Err()
	}
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		return err
	}
	return err
}
