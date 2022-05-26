package interceptors

import (
	"context"
	"strings"

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
		var (
			invalidAuthorizationToken       *epb.ErrorInfo
			invalidApp                      *epb.ErrorInfo
			invalidFirebaseCloudMessagingId *epb.ErrorInfo
			invalidDeviceId                 *epb.ErrorInfo
			invalidAccessToken              *epb.ErrorInfo
			invalidAppVersion               *epb.ErrorInfo
			invalidPlatform                 *epb.ErrorInfo
			invalidSystemVersion            *epb.ErrorInfo
			invalidModel                    *epb.ErrorInfo
		)
		var invalidArgs bool
		var st = status.New(codes.Unauthenticated, "Incorrect metadata")
		md, _ := metadata.FromIncomingContext(ctx)
		if len(md.Get("App")) == 0 {
			invalidArgs = true
			invalidApp = &epb.ErrorInfo{
				Reason: "app metadata missing",
			}
		}
		if len(md.Get("Authorization")) != 0 {
			splitValue := strings.Split(md.Get("Authorization")[0], " ")
			if splitValue[0] != "Bearer" {
				invalidArgs = true
				invalidAuthorizationToken = &epb.ErrorInfo{
					Reason: "authorization token is invalid",
				}
			}
		}
		if len(md.Get("Firebase-Cloud-Messaging-Id")) == 0 {
			invalidArgs = true
			invalidFirebaseCloudMessagingId = &epb.ErrorInfo{
				Reason: "firebase-cloud-messaging-id metadata missing",
			}
		}
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
		if len(md.Get("App-Version")) == 0 {
			invalidArgs = true
			invalidAppVersion = &epb.ErrorInfo{
				Reason: "app-version metadata missing",
			}
		}
		if invalidArgs {
			if invalidDeviceId != nil {
				st, _ = st.WithDetails(
					invalidDeviceId,
				)
			}
			if invalidAuthorizationToken != nil {
				st, _ = st.WithDetails(
					invalidAuthorizationToken,
				)
			}
			if invalidApp != nil {
				st, _ = st.WithDetails(
					invalidApp,
				)
			}
			if invalidAccessToken != nil {
				st, _ = st.WithDetails(
					invalidAccessToken,
				)
			}
			if invalidAppVersion != nil {
				st, _ = st.WithDetails(
					invalidAppVersion,
				)
			}
			if invalidPlatform != nil {
				st, _ = st.WithDetails(
					invalidPlatform,
				)
			}
			if invalidFirebaseCloudMessagingId != nil {
				st, _ = st.WithDetails(
					invalidFirebaseCloudMessagingId,
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
		var (
			invalidDeviceId      *epb.ErrorInfo
			invalidAccessToken   *epb.ErrorInfo
			invalidAppVersion    *epb.ErrorInfo
			invalidPlatform      *epb.ErrorInfo
			invalidSystemVersion *epb.ErrorInfo
			invalidModel         *epb.ErrorInfo
		)
		var st = status.New(codes.Unauthenticated, "Incorrect metadata")
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
		if len(md.Get("App-Version")) == 0 {
			invalidArgs = true
			invalidAppVersion = &epb.ErrorInfo{
				Reason: "app-version metadata missing",
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
			if invalidAppVersion != nil {
				st, _ = st.WithDetails(
					invalidAppVersion,
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
			return st.Err()
		}
		handlerErr := handler(srv, newWrappedStream(stream))
		if handlerErr != nil {
			return handlerErr
		}
		return handlerErr
	}
}
