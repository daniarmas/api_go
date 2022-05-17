package interceptors

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

func UnaryMetadataRequestInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (_ interface{}, err error) {
		var st = status.New(codes.Unauthenticated, "Incorrect metadata")
		var invalidDeviceId *epb.ErrorInfo
		var invalidAccessToken *epb.ErrorInfo
		var invalidPlatform *epb.ErrorInfo
		var invalidSystemVersion *epb.ErrorInfo
		var invalidModel *epb.ErrorInfo
		var invalidArgs bool
		md, _ := metadata.FromIncomingContext(ctx)
		if len(md.Get("Device-Id")) == 0 {
			invalidArgs = true
			invalidDeviceId = &epb.ErrorInfo{
				Reason: "device-id metadata missing",
			}
		}
		if len(md.Get("Access-Token")) == 0 {
			invalidArgs = true
			invalidAccessToken = &epb.ErrorInfo{
				Reason: "access-token metadata missing",
			}
		} else if md.Get("Access-Token")[0] != "O8pzXjp4QMk4cAD60dHeoOnxdVsDc9" {
			invalidArgs = true
			invalidAccessToken = &epb.ErrorInfo{
				Reason: "access-token is incorrect",
			}
		}
		if len(md.Get("Platform")) == 0 {
			invalidArgs = true
			invalidPlatform = &epb.ErrorInfo{
				Reason: "platform metadata missing",
			}
		}
		if len(md.Get("System-Version")) == 0 {
			invalidArgs = true
			invalidSystemVersion = &epb.ErrorInfo{
				Reason: "system-version metadata missing",
			}
		}
		if len(md.Get("Model")) == 0 {
			invalidArgs = true
			invalidModel = &epb.ErrorInfo{
				Reason: "model metadata missing",
			}
		}
		if invalidArgs {
			if invalidDeviceId != nil {
				st, _ = st.WithDetails(
					invalidDeviceId,
				)
			}
			if invalidAccessToken != nil {
				st, _ = st.WithDetails(
					invalidAccessToken,
				)
			}
			if invalidPlatform != nil {
				st, _ = st.WithDetails(
					invalidPlatform,
				)
			}
			if invalidSystemVersion != nil {
				st, _ = st.WithDetails(
					invalidSystemVersion,
				)
			}
			if invalidModel != nil {
				st, _ = st.WithDetails(
					invalidModel,
				)
			}
			return nil, st.Err()
		}
		m, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return m, err
	}
}

func StreamMetadataRequestInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) (err error) {
		var st = status.New(codes.Unauthenticated, "Incorrect metadata")
		var invalidDeviceId *epb.ErrorInfo
		var invalidAccessToken *epb.ErrorInfo
		var invalidPlatform *epb.ErrorInfo
		var invalidSystemVersion *epb.ErrorInfo
		var invalidModel *epb.ErrorInfo
		var invalidArgs bool
		md, _ := metadata.FromIncomingContext(stream.Context())
		if len(md.Get("Device-Id")) == 0 {
			invalidArgs = true
			invalidDeviceId = &epb.ErrorInfo{
				Reason: "device-id metadata missing",
			}
		}
		if len(md.Get("Access-Token")) == 0 {
			invalidArgs = true
			invalidAccessToken = &epb.ErrorInfo{
				Reason: "access-token metadata missing",
			}
		} else if md.Get("Access-Token")[0] != "O8pzXjp4QMk4cAD60dHeoOnxdVsDc9" {
			invalidArgs = true
			invalidAccessToken = &epb.ErrorInfo{
				Reason: "access-token is incorrect",
			}
		}
		if len(md.Get("Platform")) == 0 {
			invalidArgs = true
			invalidPlatform = &epb.ErrorInfo{
				Reason: "platform metadata missing",
			}
		}
		if len(md.Get("System-Version")) == 0 {
			invalidArgs = true
			invalidSystemVersion = &epb.ErrorInfo{
				Reason: "system-version metadata missing",
			}
		}
		if len(md.Get("Model")) == 0 {
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
			return st.Err()
		}
		handlerErr := handler(srv, newWrappedStream(stream))
		if handlerErr != nil {
			return handlerErr
		}
		return handlerErr
	}
}
