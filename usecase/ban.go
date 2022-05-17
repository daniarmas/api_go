package usecase

import (
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type BanService interface {
	GetBannedDevice(metadata *metadata.MD) (*models.BannedDevice, error)
	GetBannedUser(metadata *metadata.MD) (*models.BannedUser, error)
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
		bannedDevice, bannedDeviceErr = i.dao.NewBannedDeviceQuery().GetBannedDevice(tx, &models.BannedDevice{DeviceIdentifier: metadata.Get("deviceid")[0]})
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

func (i *banService) GetBannedUser(metadata *metadata.MD) (*models.BannedUser, error) {
	var bannedUser *models.BannedUser
	var bannedUserErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		bannedUser, bannedUserErr = i.dao.NewBannedUserQuery().GetBannedUser(tx, &models.BannedUser{UserId: *authorizationTokenRes.UserId})
		if bannedUserErr != nil && bannedUserErr.Error() != "record not found" {
			return bannedUserErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bannedUser, nil
}
