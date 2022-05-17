package utils

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ClientMetadata struct {
	Authorization    *string
	AccessToken      *string
	Platform         *string
	DeviceIdentifier *string
	SystemVersion    *string
	Model            *string
}

func GetMetadata(metadata *metadata.MD) (*ClientMetadata, error) {
	var st *status.Status
	var authorization, accessToken, platform, deviceIdentifier, systemVersion, model *string
	if len(metadata.Get("Authorization")) != 0 {
		splitValue := strings.Split(metadata.Get("Authorization")[0], " ")
		if len(splitValue) > 1 {
			authorization = &splitValue[1]
		}
	}
	if len(metadata.Get("Access-Token")) != 0 {
		value := metadata.Get("access-token")[0]
		accessToken = &value
	} else {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return nil, st.Err()
	}
	if len(metadata.Get("Platform")) != 0 {
		value := metadata.Get("Platform")[0]
		platform = &value
	} else {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return nil, st.Err()
	}
	if len(metadata.Get("Device-Id")) != 0 {
		value := metadata.Get("Device-Id")[0]
		deviceIdentifier = &value
	} else {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return nil, st.Err()
	}
	if len(metadata.Get("System-Version")) != 0 {
		value := metadata.Get("System-Version")[0]
		systemVersion = &value
	} else {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return nil, st.Err()
	}
	if len(metadata.Get("Model")) != 0 {
		value := metadata.Get("Model")[0]
		model = &value
	} else {
		st = status.New(codes.Unauthenticated, "Incorrect metadata")
		return nil, st.Err()
	}
	resMetadata := ClientMetadata{
		Authorization:    authorization,
		AccessToken:      accessToken,
		Platform:         platform,
		DeviceIdentifier: deviceIdentifier,
		SystemVersion:    systemVersion,
		Model:            model,
	}
	return &resMetadata, nil
}
