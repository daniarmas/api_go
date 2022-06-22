package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ApplicationService interface {
	CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest, md *utils.ClientMetadata) (*pb.Application, error)
}

type applicationService struct {
	dao repository.DAO
}

func NewApplicationService(dao repository.DAO) ApplicationService {
	return &applicationService{dao: dao}
}

func (i *applicationService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest, md *utils.ClientMetadata) (*pb.Application, error) {
	var res pb.Application
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{UserId: authorizationTokenRes.UserId, Name: "create_application"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		appRes, appErr := i.dao.NewApplicationRepository().CreateApplication(tx, &models.Application{Name: req.Application.Name, Version: req.Application.Version, Description: req.Application.Description, ExpirationTime: req.Application.ExpirationTime.AsTime()})
		if appErr != nil {
			return appErr
		}
		jwtAcessToken := datasource.JsonWebTokenMetadata{TokenId: appRes.ID}
		jwtAcessTokenErr := i.dao.NewJwtTokenRepository().CreateJwtAuthorizationToken(&jwtAcessToken)
		if jwtAcessTokenErr != nil {
			return jwtAcessTokenErr
		}
		res = pb.Application{
			Id:             appRes.ID.String(),
			Name:           appRes.Name,
			AccessToken:    *jwtAcessToken.Token,
			Version:        appRes.Version,
			Description:    appRes.Description,
			ExpirationTime: timestamppb.New(appRes.ExpirationTime),
			CreateTime:     timestamppb.New(appRes.CreateTime),
			UpdateTime:     timestamppb.New(appRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
