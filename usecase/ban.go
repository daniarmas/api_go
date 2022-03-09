package usecase

import (
	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type BanService interface {
	GetBannedDevice(metadata *metadata.MD) (*models.BannedDevice, error)
}

type banService struct {
	dao repository.DAO
}

func NewBanService(dao repository.DAO) BanService {
	return &banService{dao: dao}
}

func (i *banService) GetBannedDevice(metadata *metadata.MD) (*models.BannedDevice, error) {
	var bannedDevice *models.BannedDevice
	var bannedDeviceErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		bannedDevice, bannedDeviceErr = i.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceId: metadata.Get("deviceid")[0]}, &[]string{"create_time", "ban_expiration_time"})
		if bannedDeviceErr != nil && bannedDeviceErr.Error() != "record not found" {
			return bannedDeviceErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bannedDevice, nil
}
