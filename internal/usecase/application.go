package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/internal/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ApplicationService interface {
	ListApplication(ctx context.Context, req *pb.ListApplicationRequest, md *utils.ClientMetadata) (*pb.ListApplicationResponse, error)
	CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest, md *utils.ClientMetadata) (*pb.Application, error)
	DeleteApplication(ctx context.Context, req *pb.DeleteApplicationRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type applicationService struct {
	dao   repository.Repository
	sqldb *sqldb.Sql
}

func NewApplicationService(dao repository.Repository, sqldb *sqldb.Sql) ApplicationService {
	return &applicationService{dao: dao, sqldb: sqldb}
}

func (i *applicationService) ListApplication(ctx context.Context, req *pb.ListApplicationRequest, md *utils.ClientMetadata) (*pb.ListApplicationResponse, error) {
	var res pb.ListApplicationResponse
	var nextPage time.Time
	if req.NextPage == nil {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		err := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "read_application"}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		}
		apps, err := i.dao.NewApplicationRepository().ListApplication(tx, &entity.Application{}, &nextPage, nil)
		if err != nil {
			return err
		} else if len(*apps) > 10 {
			*apps = (*apps)[:len(*apps)-1]
			res.NextPage = timestamppb.New((*apps)[len(*apps)-1].CreateTime)
		} else if len(*apps) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*apps)[len(*apps)-1].CreateTime)
		}
		applications := make([]*pb.Application, 0, len(*apps))
		for _, i := range *apps {
			applications = append(applications, &pb.Application{
				Id:             i.ID.String(),
				Name:           i.Name,
				Version:        i.Version,
				Description:    i.Description,
				ExpirationTime: timestamppb.New(i.ExpirationTime),
				CreateTime:     timestamppb.New(i.CreateTime),
				UpdateTime:     timestamppb.New(i.UpdateTime),
			})
		}
		res.Applications = applications
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *applicationService) DeleteApplication(ctx context.Context, req *pb.DeleteApplicationRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	var res gp.Empty
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		err := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "delete_application"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		id := uuid.MustParse(req.Id)
		_, err = i.dao.NewApplicationRepository().DeleteApplication(tx, &entity.Application{ID: &id}, nil)
		if err != nil && err.Error() == "record not found" {
			return errors.New("application not found")
		} else if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *applicationService) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest, md *utils.ClientMetadata) (*pb.Application, error) {
	var res pb.Application
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "create_application"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		appRes, appErr := i.dao.NewApplicationRepository().CreateApplication(tx, &entity.Application{Name: req.Application.Name, Version: req.Application.Version, Description: req.Application.Description, ExpirationTime: req.Application.ExpirationTime.AsTime()})
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
