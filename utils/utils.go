package utils

import (
	"github.com/google/uuid"
)

func IsValidUUID(u *string) bool {
	_, err := uuid.Parse(*u)
	return err == nil
}

func UuidParse(value string) *uuid.UUID {
	if value != "" {
		val := uuid.MustParse(value)
		return &val
	} else {
		return nil
	}
}
