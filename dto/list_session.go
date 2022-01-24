package dto

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
)

type ListSessionResponse struct {
	Sessions       *[]models.Session
	ActualDeviceId uuid.UUID
}
