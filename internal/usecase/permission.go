package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	// gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PermissionService interface {
	ListPermission(ctx context.Context, req *pb.ListPermissionRequest, md *utils.ClientMetadata) (*pb.ListPermissionResponse, error)
	// CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest, md *utils.ClientMetadata) (*pb.Permission, error)
	// DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type permissionService struct {
	dao   repository.Repository
	sqldb *sqldb.Sql
}

func NewPermissionService(dao repository.Repository, sqldb *sqldb.Sql) PermissionService {
	return &permissionService{dao: dao, sqldb: sqldb}
}

func (i *permissionService) ListPermission(ctx context.Context, req *pb.ListPermissionRequest, md *utils.ClientMetadata) (*pb.ListPermissionResponse, error) {
	var res pb.ListPermissionResponse
	var nextPage time.Time
	if req.NextPage == nil {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "read_permission"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		}
		permissionsRes, err := i.dao.NewPermissionRepository().ListPermission(ctx, tx, &entity.Permission{}, nextPage)
		if err != nil {
			return err
		} else if len(*permissionsRes) > 10 {
			*permissionsRes = (*permissionsRes)[:len(*permissionsRes)-1]
			res.NextPage = timestamppb.New((*permissionsRes)[len(*permissionsRes)-1].CreateTime)
		} else if len(*permissionsRes) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*permissionsRes)[len(*permissionsRes)-1].CreateTime)
		}
		permissions := make([]*pb.Permission, 0, len(*permissionsRes))
		for _, i := range *permissionsRes {
			permissions = append(permissions, &pb.Permission{
				Id:         i.ID.String(),
				Name:       i.Name,
				CreateTime: timestamppb.New(i.CreateTime),
				UpdateTime: timestamppb.New(i.UpdateTime),
			})
		}
		res.Permissions = permissions
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// func (i *permissionService) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
// 	var res gp.Empty
// 	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
// 		_, err := i.dao.NewPermissionRepository().CheckPermission(ctx, tx, *md.AccessToken)
// 		if err != nil {
// 			return err
// 		}
// 		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
// 		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
// 		if err != nil {
// 			switch err.Error() {
// 			case "Token is expired":
// 				return errors.New("authorization token expired")
// 			case "signature is invalid":
// 				return errors.New("authorization token signature is invalid")
// 			case "token contains an invalid number of segments":
// 				return errors.New("authorization token contains an invalid number of segments")
// 			default:
// 				return err
// 			}
// 		}
// 		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
// 		if err != nil && err.Error() == "record not found" {
// 			return errors.New("unauthenticated")
// 		} else if err != nil {
// 			return err
// 		}
// 		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "delete_Permission"})
// 		if permissionErr != nil && permissionErr.Error() == "record not found" {
// 			return errors.New("permission denied")
// 		}
// 		id := uuid.MustParse(req.Id)
// 		_, err = i.dao.NewPermissionRepository().DeletePermission(ctx, tx, &entity.Permission{ID: &id}, nil)
// 		if err != nil && err.Error() == "record not found" {
// 			return errors.New("Permission not found")
// 		} else if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func (i *permissionService) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest, md *utils.ClientMetadata) (*pb.Permission, error) {
// 	var res pb.Permission
// 	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
// 		_, err := i.dao.NewPermissionRepository().CheckPermission(ctx, tx, *md.AccessToken)
// 		if err != nil {
// 			return err
// 		}
// 		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
// 		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
// 		if authorizationTokenParseErr != nil {
// 			switch authorizationTokenParseErr.Error() {
// 			case "Token is expired":
// 				return errors.New("authorization token expired")
// 			case "signature is invalid":
// 				return errors.New("authorization token signature is invalid")
// 			case "token contains an invalid number of segments":
// 				return errors.New("authorization token contains an invalid number of segments")
// 			default:
// 				return authorizationTokenParseErr
// 			}
// 		}
// 		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
// 		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
// 			return errors.New("unauthenticated")
// 		} else if authorizationTokenErr != nil {
// 			return authorizationTokenErr
// 		}
// 		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "create_Permission"})
// 		if permissionErr != nil && permissionErr.Error() == "record not found" {
// 			return errors.New("permission denied")
// 		}
// 		appRes, appErr := i.dao.NewPermissionRepository().CreatePermission(ctx, tx, &entity.Permission{Name: req.Permission.Name, Version: req.Permission.Version, Description: req.Permission.Description, ExpirationTime: req.Permission.ExpirationTime.AsTime()})
// 		if appErr != nil {
// 			return appErr
// 		}
// 		jwtAcessToken := datasource.JsonWebTokenMetadata{TokenId: appRes.ID}
// 		expirationTime := req.Permission.ExpirationTime.AsTime()
// 		jwtAcessTokenErr := i.dao.NewJwtTokenRepository().CreateJwtAccessToken(&jwtAcessToken, &expirationTime)
// 		if jwtAcessTokenErr != nil {
// 			return jwtAcessTokenErr
// 		}
// 		res = pb.Permission{
// 			Id:             appRes.ID.String(),
// 			Name:           appRes.Name,
// 			AccessToken:    *jwtAcessToken.Token,
// 			Version:        appRes.Version,
// 			Description:    appRes.Description,
// 			ExpirationTime: timestamppb.New(appRes.ExpirationTime),
// 			CreateTime:     timestamppb.New(appRes.CreateTime),
// 			UpdateTime:     timestamppb.New(appRes.UpdateTime),
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }
