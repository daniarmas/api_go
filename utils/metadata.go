package utils

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type ClientMetadata struct {
	Authorization            *string
	FirebaseCloudMessagingId *string
	AccessToken              *string
	Platform                 *string
	DeviceIdentifier         *string
	SystemVersion            *string
	AppVersion               *string
	Model                    *string
	App                      *string
}

func GetMetadata(ctx context.Context) *ClientMetadata {
	var authorization *string
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md.Get("Authorization")) != 0 {
		splitValue := strings.Split(md.Get("Authorization")[0], " ")
		if len(splitValue) > 1 {
			authorization = &splitValue[1]
		}
	}
	app := md.Get("App")[0]
	appVersion := md.Get("App-Version")[0]
	firebaseCloudMessagingId := md.Get("Firebase-Cloud-Messaging-Id")[0]
	accessToken := md.Get("Access-Token")[0]
	platform := md.Get("Platform")[0]
	deviceIdentifier := md.Get("Device-Id")[0]
	systemVersion := md.Get("System-Version")[0]
	model := md.Get("Model")[0]
	resMetadata := ClientMetadata{
		App:                      &app,
		Authorization:            authorization,
		FirebaseCloudMessagingId: &firebaseCloudMessagingId,
		AppVersion:               &appVersion,
		AccessToken:              &accessToken,
		Platform:                 &platform,
		DeviceIdentifier:         &deviceIdentifier,
		SystemVersion:            &systemVersion,
		Model:                    &model,
	}
	return &resMetadata
}
