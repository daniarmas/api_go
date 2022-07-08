package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/internal/repository"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type BanService interface {
	GetBannedDevice(metadata *metadata.MD) (*entity.BannedDevice, error)
	GetBannedUser(metadata *metadata.MD) (*entity.BannedUser, error)
}

type banService struct {
	dao   repository.Repository
	sqldb *sqldb.Sql
}

func NewBanService(dao repository.Repository, sqldb *sqldb.Sql) BanService {
	return &banService{dao: dao, sqldb: sqldb}
}

func (i *banService) GetBannedDevice(metadata *metadata.MD) (*entity.BannedDevice, error) {
	var bannedDevice *entity.BannedDevice
	var bannedDeviceErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		bannedDevice, bannedDeviceErr = i.dao.NewBannedDeviceRepository().GetBannedDevice(tx, &entity.BannedDevice{DeviceIdentifier: metadata.Get("deviceid")[0]}, nil)
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

func (i *banService) GetBannedUser(metadata *metadata.MD) (*entity.BannedUser, error) {
	var ctx = context.Background()
	var bannedUser *entity.BannedUser
	var bannedUserErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		bannedUser, bannedUserErr = i.dao.NewBannedUserRepository().GetBannedUser(tx, &entity.BannedUser{UserId: authorizationTokenRes.UserId}, nil)
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
