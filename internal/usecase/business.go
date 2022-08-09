package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type BusinessService interface {
	Feed(ctx context.Context, req *pb.FeedRequest, meta *utils.ClientMetadata) (*pb.FeedResponse, error)
	BusinessIsInRange(ctx context.Context, req *pb.BusinessIsInRangeRequest, meta *utils.ClientMetadata) error
	GetBusiness(ctx context.Context, req *pb.GetBusinessRequest, meta *utils.ClientMetadata) (*pb.Business, error)
	GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest, md *utils.ClientMetadata) (*pb.Business, error)
	CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest, md *utils.ClientMetadata) (*pb.CreateBusinessResponse, error)
	UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest, md *utils.ClientMetadata) (*pb.Business, error)
	CreatePartnerApplication(ctx context.Context, req *pb.CreatePartnerApplicationRequest, md *utils.ClientMetadata) (*pb.PartnerApplication, error)
	ListPartnerApplication(ctx context.Context, req *pb.ListPartnerApplicationRequest, md *utils.ClientMetadata) (*pb.ListPartnerApplicationResponse, error)
	UpdatePartnerApplication(ctx context.Context, req *pb.UpdatePartnerApplicationRequest, md *utils.ClientMetadata) (*pb.PartnerApplication, error)
	ListBusinessRole(ctx context.Context, req *pb.ListBusinessRoleRequest, md *utils.ClientMetadata) (*pb.ListBusinessRoleResponse, error)
	CreateBusinessRole(ctx context.Context, req *pb.CreateBusinessRoleRequest, md *utils.ClientMetadata) (*pb.BusinessRole, error)
	UpdateBusinessRole(ctx context.Context, req *pb.UpdateBusinessRoleRequest, md *utils.ClientMetadata) (*pb.BusinessRole, error)
	DeleteBusinessRole(ctx context.Context, req *pb.DeleteBusinessRoleRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	ModifyBusinessRolePermission(ctx context.Context, req *pb.ModifyBusinessRolePermissionRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type businessService struct {
	config *config.Config
	dao    repository.Repository
	stDb   *sql.DB
	sqldb  *sqldb.Sql
}

func NewBusinessService(dao repository.Repository, config *config.Config, stDb *sql.DB, sqldb *sqldb.Sql) BusinessService {
	return &businessService{dao: dao, config: config, stDb: stDb, sqldb: sqldb}
}

func (i *businessService) BusinessIsInRange(ctx context.Context, req *pb.BusinessIsInRangeRequest, meta *utils.ClientMetadata) error {
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *meta.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: meta.Authorization}
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
		_, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		businessId := uuid.MustParse(req.BusinessId)
		businessIsInRange, err := i.dao.NewBusinessRepository().BusinessIsInRange(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)}, &businessId)
		if err != nil {
			return err
		}
		if !*businessIsInRange {
			return errors.New("business is not in range")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *businessService) ModifyBusinessRolePermission(ctx context.Context, req *pb.ModifyBusinessRolePermissionRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	var res gp.Empty
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		businessRoleId := uuid.MustParse(req.BusinessRoleId)
		businessRoleRes, err := i.dao.NewBusinessRoleRepository().GetBusinessRole(tx, &entity.BusinessRole{ID: &businessRoleId}, nil)
		if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_role", BusinessId: businessRoleRes.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		}
		unionBusinessRoleAndPermission := make([]entity.UnionBusinessRoleAndPermission, 0, len(req.PermissionIds))
		permissionIds := make([]uuid.UUID, 0, len(req.PermissionIds))
		for _, i := range req.PermissionIds {
			permissionId := uuid.MustParse(i)
			unionBusinessRoleAndPermission = append(unionBusinessRoleAndPermission, entity.UnionBusinessRoleAndPermission{
				BusinessRoleId: &businessRoleId,
				PermissionId:   &permissionId,
			})
			permissionIds = append(permissionIds, permissionId)
		}
		permissionsRes, err := i.dao.NewPermissionRepository().ListPermissionByIdAll(tx, &entity.Permission{}, &permissionIds)
		if err != nil {
			return err
		}
		unionBusinessRoleAndUser, err := i.dao.NewUnionBusinessRoleAndUserRepository().ListUnionBusinessRoleAndUserAll(tx, &entity.UnionBusinessRoleAndUser{BusinessRoleId: &businessRoleId})
		if err != nil {
			return err
		}
		userPermissions := make([]entity.UserPermission, 0, len(req.PermissionIds))
		for _, c := range *unionBusinessRoleAndUser {
			for _, i := range permissionIds {
				for _, y := range *permissionsRes {
					if *y.ID == i {
						userPermissions = append(userPermissions, entity.UserPermission{
							Name:           y.Name,
							UserId:         c.UserId,
							BusinessId:     businessRoleRes.BusinessId,
							BusinessRoleId: businessRoleRes.ID,
							PermissionId:   y.ID,
						})
					}
				}
			}
		}
		_, err = i.dao.NewUnionBusinessRoleAndPermissionRepository().DeleteUnionBusinessRoleAndPermission(tx, &entity.UnionBusinessRoleAndPermission{BusinessRoleId: &businessRoleId}, nil)
		if err != nil && err.Error() == "record not found" {
			return errors.New("business role not found")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().DeleteUserPermissionByBusinessRoleId(tx, &entity.UserPermission{BusinessRoleId: &businessRoleId})
		if err != nil {
			return err
		}
		_, err = i.dao.NewUnionBusinessRoleAndPermissionRepository().CreateUnionBusinessRoleAndPermission(tx, &unionBusinessRoleAndPermission)
		if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().CreateUserPermission(tx, &userPermissions)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) UpdateBusinessRole(ctx context.Context, req *pb.UpdateBusinessRoleRequest, md *utils.ClientMetadata) (*pb.BusinessRole, error) {
	var res pb.BusinessRole
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		businessRolesRes, err := i.dao.NewBusinessRoleRepository().UpdateBusinessRole(tx, &entity.BusinessRole{ID: &id}, &entity.BusinessRole{Name: req.BusinessRole.Name})
		if err != nil && err.Error() == "record not found" {
			return errors.New("business role not found")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_role", BusinessId: businessRolesRes.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		}
		res = pb.BusinessRole{
			Id:         businessRolesRes.ID.String(),
			Name:       businessRolesRes.Name,
			BusinessId: businessRolesRes.BusinessId.String(),
			CreateTime: timestamppb.New(businessRolesRes.CreateTime),
			UpdateTime: timestamppb.New(businessRolesRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) DeleteBusinessRole(ctx context.Context, req *pb.DeleteBusinessRoleRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		id := uuid.MustParse(req.Id)
		businessRolesRes, err := i.dao.NewBusinessRoleRepository().DeleteBusinessRole(tx, &entity.BusinessRole{ID: &id}, nil)
		if err != nil {
			return err
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "delete_role", BusinessId: (*businessRolesRes)[0].BusinessId}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		unionBusinessRoleAndPermRes, err := i.dao.NewUnionBusinessRoleAndPermissionRepository().DeleteUnionBusinessRoleAndPermission(tx, &entity.UnionBusinessRoleAndPermission{BusinessRoleId: (*businessRolesRes)[0].ID}, nil)
		if err != nil {
			return err
		}
		userPermissionIds := make([]uuid.UUID, 0, len(*unionBusinessRoleAndPermRes))
		for _, i := range *unionBusinessRoleAndPermRes {
			userPermissionIds = append(userPermissionIds, *i.PermissionId)
		}
		_, err = i.dao.NewUserPermissionRepository().DeleteUserPermissionByPermissionId(tx, &userPermissionIds)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
}

func (i *businessService) CreateBusinessRole(ctx context.Context, req *pb.CreateBusinessRoleRequest, md *utils.ClientMetadata) (*pb.BusinessRole, error) {
	var res pb.BusinessRole
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		businessId := uuid.MustParse(req.BusinessRole.BusinessId)
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "create_role", BusinessId: &businessId}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		businessRolesRes, err := i.dao.NewBusinessRoleRepository().CreateBusinessRole(tx, &entity.BusinessRole{Name: req.BusinessRole.Name, BusinessId: &businessId})
		if err != nil {
			return err
		}
		businessRolePermissions := make([]entity.UnionBusinessRoleAndPermission, 0, len(req.BusinessRole.Permissions))
		for _, i := range req.BusinessRole.Permissions {
			permissionId := uuid.MustParse(i.Id)
			businessRolePermissions = append(businessRolePermissions, entity.UnionBusinessRoleAndPermission{
				PermissionId:   &permissionId,
				BusinessRoleId: businessRolesRes.ID,
			})
		}
		_, err = i.dao.NewUnionBusinessRoleAndPermissionRepository().CreateUnionBusinessRoleAndPermission(tx, &businessRolePermissions)
		if err != nil {
			return err
		}
		res = pb.BusinessRole{
			Id:         businessRolesRes.ID.String(),
			Name:       businessRolesRes.Name,
			BusinessId: businessRolesRes.BusinessId.String(),
			CreateTime: timestamppb.New(businessRolesRes.CreateTime),
			UpdateTime: timestamppb.New(businessRolesRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) ListBusinessRole(ctx context.Context, req *pb.ListBusinessRoleRequest, md *utils.ClientMetadata) (*pb.ListBusinessRoleResponse, error) {
	var res pb.ListBusinessRoleResponse
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
		businessId := uuid.MustParse(req.BusinessId)
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "read_role", BusinessId: &businessId}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		businessRolesRes, err := i.dao.NewBusinessRoleRepository().ListBusinessRole(tx, &entity.BusinessRole{}, &nextPage, nil)
		if err != nil {
			return err
		} else if len(*businessRolesRes) > 10 {
			*businessRolesRes = (*businessRolesRes)[:len(*businessRolesRes)-1]
			res.NextPage = timestamppb.New((*businessRolesRes)[len(*businessRolesRes)-1].CreateTime)
		} else if len(*businessRolesRes) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*businessRolesRes)[len(*businessRolesRes)-1].CreateTime)
		}
		businessRoles := make([]*pb.BusinessRole, 0, len(*businessRolesRes))
		for _, i := range *businessRolesRes {
			businessRoles = append(businessRoles, &pb.BusinessRole{
				Id:         i.ID.String(),
				Name:       i.Name,
				BusinessId: i.BusinessId.String(),
				CreateTime: timestamppb.New(i.CreateTime),
				UpdateTime: timestamppb.New(i.UpdateTime),
			})
		}
		res.BusinessRoles = businessRoles
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) UpdatePartnerApplication(ctx context.Context, req *pb.UpdatePartnerApplicationRequest, md *utils.ClientMetadata) (*pb.PartnerApplication, error) {
	var res pb.PartnerApplication
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		id := uuid.MustParse(req.Id)
		getPartnerAppRes, getPartnerAppErr := i.dao.NewPartnerApplicationRepository().GetPartnerApplication(tx, &entity.PartnerApplication{ID: &id}, &[]string{"user_id"})
		if getPartnerAppErr != nil && getPartnerAppErr.Error() == "record not found" {
			return errors.New("partner application not found")
		} else if getPartnerAppErr != nil {
			return getPartnerAppErr
		}
		if req.PartnerApplication.Status == pb.PartnerApplicationStatus_PartnerApplicationStatusApproved || req.PartnerApplication.Status == pb.PartnerApplicationStatus_PartnerApplicationStatusRejected {
			_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_partner_application"}, &[]string{"id"})
			if permissionErr != nil && permissionErr.Error() == "record not found" {
				return errors.New("permission denied")
			}
			if req.PartnerApplication.Status == pb.PartnerApplicationStatus_PartnerApplicationStatusApproved {
				userId := uuid.MustParse(req.PartnerApplication.UserId)
				_, businessUserErr := i.dao.NewBusinessUserRepository().CreateBusinessUser(tx, &entity.BusinessUser{IsBusinessOwner: true, UserId: &userId})
				if businessUserErr != nil {
					return businessUserErr
				}
			}
		} else if req.PartnerApplication.Status == pb.PartnerApplicationStatus_PartnerApplicationStatusCanceled {
			if *getPartnerAppRes.UserId != *authorizationTokenRes.UserId {
				return errors.New("permission denied")
			}
		}
		updatePartnerAppRes, updatePartnerAppErr := i.dao.NewPartnerApplicationRepository().UpdatePartnerApplication(tx, &entity.PartnerApplication{ID: &id}, &entity.PartnerApplication{Status: req.PartnerApplication.Status.String()})
		if updatePartnerAppErr != nil {
			return updatePartnerAppErr
		}
		res = pb.PartnerApplication{
			Id:             updatePartnerAppRes.ID.String(),
			BusinessName:   updatePartnerAppRes.BusinessName,
			Coordinates:    &pb.Point{Latitude: updatePartnerAppRes.Coordinates.FlatCoords()[0], Longitude: updatePartnerAppRes.Coordinates.FlatCoords()[1]},
			Description:    updatePartnerAppRes.Description,
			UserId:         updatePartnerAppRes.UserId.String(),
			MunicipalityId: updatePartnerAppRes.MunicipalityId.String(),
			ProvinceId:     updatePartnerAppRes.ProvinceId.String(),
			Status:         *utils.ParsePartnerApplicationStatus(&updatePartnerAppRes.Status),
			CreateTime:     timestamppb.New(updatePartnerAppRes.CreateTime),
			UpdateTime:     timestamppb.New(updatePartnerAppRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) ListPartnerApplication(ctx context.Context, req *pb.ListPartnerApplicationRequest, md *utils.ClientMetadata) (*pb.ListPartnerApplicationResponse, error) {
	var res pb.ListPartnerApplicationResponse
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
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "read_partner_application"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("not permission")
		}
		partnerApplicationsRes, partnerApplicationsErr := i.dao.NewPartnerApplicationRepository().ListPartnerApplication(tx, nil, &nextPage, nil)
		if partnerApplicationsErr != nil {
			return partnerApplicationsErr
		} else if len(*partnerApplicationsRes) > 10 {
			*partnerApplicationsRes = (*partnerApplicationsRes)[:len(*partnerApplicationsRes)-1]
			res.NextPage = timestamppb.New((*partnerApplicationsRes)[len(*partnerApplicationsRes)-1].CreateTime)
		} else if len(*partnerApplicationsRes) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*partnerApplicationsRes)[len(*partnerApplicationsRes)-1].CreateTime)
		}
		partnerApplications := make([]*pb.PartnerApplication, 0, len(*partnerApplicationsRes))
		for _, i := range *partnerApplicationsRes {
			partnerApplications = append(partnerApplications, &pb.PartnerApplication{
				Id:             i.ID.String(),
				BusinessName:   i.BusinessName,
				Coordinates:    &pb.Point{Latitude: i.Coordinates.FlatCoords()[0], Longitude: i.Coordinates.FlatCoords()[1]},
				Description:    i.Description,
				Status:         *utils.ParsePartnerApplicationStatus(&i.Status),
				UserId:         i.UserId.String(),
				MunicipalityId: i.MunicipalityId.String(),
				ProvinceId:     i.ProvinceId.String(),
				CreateTime:     timestamppb.New(i.CreateTime),
				UpdateTime:     timestamppb.New(i.UpdateTime),
			})
		}
		res.PartnerApplications = partnerApplications
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) CreatePartnerApplication(ctx context.Context, req *pb.CreatePartnerApplicationRequest, md *utils.ClientMetadata) (*pb.PartnerApplication, error) {
	var res pb.PartnerApplication
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		businessUserRes, businessUserErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &entity.BusinessUser{UserId: authorizationTokenRes.UserId}, &[]string{"id"})
		if businessUserErr != nil && businessUserErr.Error() != "record not found" {
			return businessUserErr
		}
		if businessUserRes != nil {
			return errors.New("already register as business user")
		}
		businessRes, businessErr := i.dao.NewBusinessRepository().GetBusiness(tx, &entity.Business{Name: req.PartnerApplication.BusinessName}, &[]string{"id"})
		if businessErr != nil && businessErr.Error() != "record not found" {
			return businessErr
		}
		if businessRes != nil {
			return errors.New("already exists a business with that name")
		}
		municipalityId := uuid.MustParse(req.PartnerApplication.MunicipalityId)
		provinceId := uuid.MustParse(req.PartnerApplication.ProvinceId)
		location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.PartnerApplication.Coordinates.Latitude, req.PartnerApplication.Coordinates.Longitude}).SetSRID(4326)}
		createPartnerAppRes, createPartnerAppErr := i.dao.NewPartnerApplicationRepository().CreatePartnerApplication(tx, &entity.PartnerApplication{BusinessName: req.PartnerApplication.BusinessName, Description: req.PartnerApplication.Description, ProvinceId: &provinceId, MunicipalityId: &municipalityId, UserId: authorizationTokenRes.UserId, Coordinates: location})
		if createPartnerAppErr != nil {
			return createPartnerAppErr
		}
		res = pb.PartnerApplication{
			Id:             createPartnerAppRes.ID.String(),
			BusinessName:   createPartnerAppRes.BusinessName,
			Coordinates:    &pb.Point{Latitude: createPartnerAppRes.Coordinates.FlatCoords()[0], Longitude: createPartnerAppRes.Coordinates.FlatCoords()[1]},
			Description:    createPartnerAppRes.Description,
			UserId:         createPartnerAppRes.UserId.String(),
			Status:         *utils.ParsePartnerApplicationStatus(&createPartnerAppRes.Status),
			MunicipalityId: createPartnerAppRes.MunicipalityId.String(),
			ProvinceId:     createPartnerAppRes.ProvinceId.String(),
			CreateTime:     timestamppb.New(createPartnerAppRes.CreateTime),
			UpdateTime:     timestamppb.New(createPartnerAppRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *businessService) UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest, md *utils.ClientMetadata) (*pb.Business, error) {
	var businessRes *entity.Business
	var businessErr error
	id := uuid.MustParse(req.Id)
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &entity.BusinessUser{UserId: authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		businessIsOpenRes, businessIsOpenErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &entity.BusinessSchedule{BusinessId: &id}, "OrderTypePickUp")
		if businessIsOpenErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessIsOpenErr
		} else if businessIsOpenRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &entity.BusinessSchedule{BusinessId: &id}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{BusinessId: &id}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		getBusinessRes, getBusinessErr := i.dao.NewBusinessRepository().GetBusiness(tx, &entity.Business{ID: &id}, nil)
		if getBusinessErr != nil {
			return getBusinessErr
		}
		if req.HighQualityPhoto != "" || req.LowQualityPhoto != "" || req.Thumbnail != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.HighQualityPhoto)
			if hqErr != nil && hqErr.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.LowQualityPhoto)
			if lqErr != nil && lqErr.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.Thumbnail)
			if tnErr != nil && tnErr.Error() == "ObjectMissing" {
				return errors.New("ThumbnailObject missing")
			} else if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getBusinessRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.BusinessAvatarBulkName, Object: getBusinessRes.HighQualityPhoto})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getBusinessRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.BusinessAvatarBulkName, Object: getBusinessRes.LowQualityPhoto})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getBusinessRes.Thumbnail}, minio.CopySrcOptions{Bucket: i.config.BusinessAvatarBulkName, Object: getBusinessRes.Thumbnail})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.BusinessAvatarBulkName, getBusinessRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.BusinessAvatarBulkName, getBusinessRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.BusinessAvatarBulkName, getBusinessRes.Thumbnail, minio.RemoveObjectOptions{})
			if rmThErr != nil {
				return rmThErr
			}
		}
		var provinceId uuid.UUID
		var municipalityId uuid.UUID
		if req.ProvinceId != "" {
			provinceId = uuid.MustParse(req.ProvinceId)
		}
		if req.MunicipalityId != "" {
			municipalityId = uuid.MustParse(req.MunicipalityId)
		}
		businessRes, businessErr = i.dao.NewBusinessRepository().UpdateBusiness(tx, &entity.Business{
			Name:                  req.Name,
			Address:               req.Address,
			HighQualityPhoto:      i.config.BusinessAvatarBulkName + "/" + req.HighQualityPhoto,
			LowQualityPhoto:       i.config.BusinessAvatarBulkName + "/" + req.LowQualityPhoto,
			Thumbnail:             i.config.BusinessAvatarBulkName + "/" + req.Thumbnail,
			BlurHash:              req.BlurHash,
			TimeMarginOrderMonth:  req.TimeMarginOrderMonth,
			TimeMarginOrderDay:    req.TimeMarginOrderDay,
			TimeMarginOrderHour:   req.TimeMarginOrderHour,
			TimeMarginOrderMinute: req.TimeMarginOrderMinute,
			DeliveryPriceCup:      req.DeliveryPriceCup,
			ToPickUp:              req.ToPickUp,
			HomeDelivery:          req.HomeDelivery,
			ProvinceId:            &provinceId,
			MunicipalityId:        &municipalityId,
		}, &entity.Business{ID: &id})
		if businessErr != nil {
			return businessErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, DeliveryPriceCup: businessRes.DeliveryPriceCup, TimeMarginOrderMonth: businessRes.TimeMarginOrderMonth, TimeMarginOrderDay: businessRes.TimeMarginOrderDay, TimeMarginOrderHour: businessRes.TimeMarginOrderHour, TimeMarginOrderMinute: businessRes.TimeMarginOrderMinute, ToPickUp: businessRes.ToPickUp, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), CreateTime: timestamppb.New(businessRes.CreateTime), UpdateTime: timestamppb.New(businessRes.UpdateTime)}, nil
}

func (i *businessService) CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest, md *utils.ClientMetadata) (*pb.CreateBusinessResponse, error) {
	var businessRes *entity.Business
	var businessErr error
	var response pb.CreateBusinessResponse
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
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
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &entity.BusinessUser{UserId: authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("HighQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("LowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.BusinessAvatarBulkName, req.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("ThumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		provinceId := uuid.MustParse(req.ProvinceId)
		municipalityId := uuid.MustParse(req.MunicipalityId)
		businessBrandId := uuid.MustParse(req.BusinessBrandId)
		businessRes, businessErr = i.dao.NewBusinessRepository().CreateBusiness(tx, &entity.Business{
			Name:                  req.Name,
			Address:               req.Address,
			HighQualityPhoto:      req.HighQualityPhoto,
			LowQualityPhoto:       req.LowQualityPhoto,
			Thumbnail:             req.Thumbnail,
			BlurHash:              req.BlurHash,
			TimeMarginOrderMonth:  req.TimeMarginOrderMonth,
			TimeMarginOrderDay:    req.TimeMarginOrderDay,
			TimeMarginOrderHour:   req.TimeMarginOrderHour,
			TimeMarginOrderMinute: req.TimeMarginOrderMinute,
			DeliveryPriceCup:      req.DeliveryPriceCup,
			ToPickUp:              req.ToPickUp,
			HomeDelivery:          req.HomeDelivery,
			Coordinates:           ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Coordinates.Latitude, req.Coordinates.Longitude}).SetSRID(4326)},
			ProvinceId:            &provinceId,
			MunicipalityId:        &municipalityId,
			BusinessBrandId:       &businessBrandId,
		})
		if businessErr != nil {
			return businessErr
		}
		response.Business = &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, DeliveryPriceCup: businessRes.DeliveryPriceCup, TimeMarginOrderMonth: businessRes.TimeMarginOrderMonth, TimeMarginOrderDay: businessRes.TimeMarginOrderDay, TimeMarginOrderHour: businessRes.TimeMarginOrderHour, TimeMarginOrderMinute: businessRes.TimeMarginOrderMinute, ToPickUp: businessRes.ToPickUp, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), CreateTime: timestamppb.New(businessRes.CreateTime), UpdateTime: timestamppb.New(businessRes.UpdateTime), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.FlatCoords()[0], Longitude: businessRes.Coordinates.FlatCoords()[1]}}
		var unionBusinessAndMunicipalities = make([]*entity.UnionBusinessAndMunicipality, 0, len(req.Municipalities))
		for _, item := range req.Municipalities {
			municipalityId := uuid.MustParse(item)
			unionBusinessAndMunicipalities = append(unionBusinessAndMunicipalities, &entity.UnionBusinessAndMunicipality{
				BusinessId:     businessRes.ID,
				MunicipalityId: &municipalityId,
			})
		}
		unionBusinessAndMunicipalityRes, unionBusinessAndMunicipalityErr := i.dao.NewUnionBusinessAndMunicipalityRepository().BatchCreateUnionBusinessAndMunicipality(tx, unionBusinessAndMunicipalities)
		if unionBusinessAndMunicipalityErr != nil {
			return unionBusinessAndMunicipalityErr
		}
		unionBusinessAndMunicipalityIds := make([]string, 0, len(unionBusinessAndMunicipalityRes))
		for _, item := range unionBusinessAndMunicipalityRes {
			unionBusinessAndMunicipalityIds = append(unionBusinessAndMunicipalityIds, item.ID.String())
		}
		_, unionBusinessAndMunicipalityWithMunicipalityErr := i.dao.NewUnionBusinessAndMunicipalityRepository().ListUnionBusinessAndMunicipalityWithMunicipality(tx, unionBusinessAndMunicipalityIds)
		if unionBusinessAndMunicipalityWithMunicipalityErr != nil {
			return unionBusinessAndMunicipalityWithMunicipalityErr
		}
		// response.UnionBusinessAndMunicipalityWithMunicipality = unionBusinessAndMunicipalityWithMunicipalityRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *businessService) Feed(ctx context.Context, req *pb.FeedRequest, meta *utils.ClientMetadata) (*pb.FeedResponse, error) {
	var businessRes *[]entity.Business
	var businessResAdd *[]entity.Business
	var businessErr, businessErrAdd error
	var response pb.FeedResponse
	var businessCategories []*pb.BusinessCategory
	provinceId := uuid.MustParse(req.ProvinceId)
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *meta.AccessToken)
		if err != nil {
			return err
		}
		if req.SearchMunicipalityType == pb.SearchMunicipalityType_More {
			businessRes, businessErr = v.dao.NewBusinessRepository().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, req.NextPage, false, req.HomeDelivery, req.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			if len(*businessRes) > 5 {
				*businessRes = (*businessRes)[:len(*businessRes)-1]
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
				response.SearchMunicipalityType = pb.SearchMunicipalityType_More
			} else if len(*businessRes) <= 5 && len(*businessRes) != 0 {
				length := 5 - len(*businessRes)
				businessResAdd, businessErrAdd = v.dao.NewBusinessRepository().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, int32(length), req.ProvinceId, req.MunicipalityId, 0, true, req.HomeDelivery, req.ToPickUp)
				if businessErrAdd != nil {
					return businessErrAdd
				}
				if businessResAdd != nil {
					if len(*businessResAdd) > length {
						*businessResAdd = (*businessResAdd)[:len(*businessResAdd)-1]
					}
					*businessRes = append(*businessRes, *businessResAdd...)
				}
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
				response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore
			} else if len(*businessRes) == 0 {
				businessRes, businessErr = v.dao.NewBusinessRepository().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, 0, true, req.HomeDelivery, req.ToPickUp)
				if businessErr != nil {
					return businessErr
				}
				if len(*businessRes) > 5 {
					*businessRes = (*businessRes)[:len(*businessRes)-1]
					response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
					response.SearchMunicipalityType = pb.SearchMunicipalityType_More
				}
			}
		} else {
			businessRes, businessErr = v.dao.NewBusinessRepository().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, req.NextPage, true, req.HomeDelivery, req.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			if businessRes != nil && len(*businessRes) > 5 {
				*businessRes = (*businessRes)[:len(*businessRes)-1]
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			} else if businessRes != nil && len(*businessRes) <= 5 && len(*businessRes) != 0 {
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			} else {
				response.NextPage = req.NextPage
			}
			response.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore
		}
		if businessRes != nil {
			businessResponse := make([]*pb.Business, 0, len(*businessRes))
			for _, e := range *businessRes {
				businessResponse = append(businessResponse, &pb.Business{
					Id:                    e.ID.String(),
					Name:                  e.Name,
					HighQualityPhoto:      e.HighQualityPhoto,
					HighQualityPhotoUrl:   v.config.BusinessAvatarBulkName + "/" + e.HighQualityPhoto,
					LowQualityPhoto:       e.LowQualityPhoto,
					LowQualityPhotoUrl:    v.config.BusinessAvatarBulkName + "/" + e.LowQualityPhoto,
					Thumbnail:             e.Thumbnail,
					ThumbnailUrl:          v.config.BusinessAvatarBulkName + "/" + e.Thumbnail,
					BlurHash:              e.BlurHash,
					Address:               e.Address,
					DeliveryPriceCup:      e.DeliveryPriceCup,
					TimeMarginOrderMonth:  e.TimeMarginOrderMonth,
					TimeMarginOrderDay:    e.TimeMarginOrderDay,
					TimeMarginOrderHour:   e.TimeMarginOrderHour,
					TimeMarginOrderMinute: e.TimeMarginOrderMinute,
					ToPickUp:              e.ToPickUp,
					HomeDelivery:          e.HomeDelivery,
					BusinessBrandId:       e.BusinessBrandId.String(),
					ProvinceId:            e.ProvinceId.String(),
					MunicipalityId:        e.MunicipalityId.String(),
					Cursor:                int32(e.Cursor),
				})
			}
			response.Businesses = businessResponse
		}
		if req.Categories {
			res, err := v.dao.NewBusinessCategoryRepository().ListBusinessCategory(tx, &entity.BusinessCategory{ProvinceId: &provinceId}, nil)
			if err != nil {
				return err
			}
			for _, i := range *res {
				businessCategories = append(businessCategories, &pb.BusinessCategory{
					Id:             i.ID.String(),
					Name:           i.Name,
					ProvinceId:     i.ProvinceId.String(),
					MunicipalityId: i.MunicipalityId.String(),
					CreateTime:     timestamppb.New(i.CreateTime),
					UpdateTime:     timestamppb.New(i.UpdateTime),
				})
			}
			response.Categories = businessCategories
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *businessService) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest, meta *utils.ClientMetadata) (*pb.Business, error) {
	var businessRes *entity.Business
	var businessCollectionRes *[]entity.BusinessCollection
	var businessErr, businessCollectionErr error
	var itemsCategoryResponse []*pb.BusinessCollection
	var schedule *entity.BusinessSchedule
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *meta.AccessToken)
		if err != nil {
			return err
		}
		businessId := uuid.MustParse(req.Id)
		businessRes, businessErr = v.dao.NewBusinessRepository().GetBusiness(tx, &entity.Business{ID: &businessId}, nil)
		if businessErr != nil && businessErr.Error() == "record not found" {
			return errors.New("business not found")
		} else if businessErr != nil {
			return businessErr
		}
		schedule, err = v.dao.NewBusinessScheduleRepository().GetBusinessSchedule(tx, &entity.BusinessSchedule{BusinessId: businessRes.ID}, nil)
		if err != nil {
			return err
		}
		businessCollectionRes, businessCollectionErr = v.dao.NewBusinessCollectionRepository().ListBusinessCollection(tx, &entity.BusinessCollection{BusinessId: &businessId}, nil)
		if businessCollectionErr != nil {
			return businessCollectionErr
		}
		itemsCategoryResponse = make([]*pb.BusinessCollection, 0, len(*businessCollectionRes))
		for _, e := range *businessCollectionRes {
			itemsCategoryResponse = append(itemsCategoryResponse, &pb.BusinessCollection{
				Id:         e.ID.String(),
				Name:       e.Name,
				BusinessId: e.BusinessId.String(),
				Index:      e.Index,
				CreateTime: timestamppb.New(e.CreateTime),
				UpdateTime: timestamppb.New(e.UpdateTime),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
	highQualityPhotoUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.HighQualityPhoto
	lowQualityPhotoUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.LowQualityPhoto
	thumbnailUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.Thumbnail
	return &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, ToPickUp: businessRes.ToPickUp, DeliveryPriceCup: businessRes.DeliveryPriceCup, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.Coords()[1], Longitude: businessRes.Coordinates.Coords()[0]}, HighQualityPhotoUrl: highQualityPhotoUrl, LowQualityPhotoUrl: lowQualityPhotoUrl, ThumbnailUrl: thumbnailUrl, BusinessCollections: itemsCategoryResponse, BusinessSchedule: &pb.BusinessSchedule{
		Id:                         schedule.ID.String(),
		FirstOpeningTimeSunday:     timestamppb.New(schedule.FirstOpeningTimeSunday),
		FirstClosingTimeSunday:     timestamppb.New(schedule.FirstClosingTimeSunday),
		FirstOpeningTimeMonday:     timestamppb.New(schedule.FirstOpeningTimeMonday),
		FirstClosingTimeMonday:     timestamppb.New(schedule.FirstClosingTimeMonday),
		FirstOpeningTimeTuesday:    timestamppb.New(schedule.FirstOpeningTimeTuesday),
		FirstClosingTimeTuesday:    timestamppb.New(schedule.FirstClosingTimeTuesday),
		FirstOpeningTimeWednesday:  timestamppb.New(schedule.FirstOpeningTimeWednesday),
		FirstClosingTimeWednesday:  timestamppb.New(schedule.FirstClosingTimeWednesday),
		FirstOpeningTimeThursday:   timestamppb.New(schedule.FirstOpeningTimeThursday),
		FirstClosingTimeThursday:   timestamppb.New(schedule.FirstClosingTimeThursday),
		FirstOpeningTimeFriday:     timestamppb.New(schedule.FirstOpeningTimeFriday),
		FirstClosingTimeFriday:     timestamppb.New(schedule.FirstClosingTimeFriday),
		FirstOpeningTimeSaturday:   timestamppb.New(schedule.FirstOpeningTimeSaturday),
		FirstClosingTimeSaturday:   timestamppb.New(schedule.FirstClosingTimeSaturday),
		SecondOpeningTimeSunday:    timestamppb.New(schedule.SecondOpeningTimeSunday),
		SecondClosingTimeSunday:    timestamppb.New(schedule.SecondClosingTimeSunday),
		SecondOpeningTimeMonday:    timestamppb.New(schedule.SecondOpeningTimeMonday),
		SecondClosingTimeMonday:    timestamppb.New(schedule.SecondClosingTimeMonday),
		SecondOpeningTimeTuesday:   timestamppb.New(schedule.SecondOpeningTimeTuesday),
		SecondClosingTimeTuesday:   timestamppb.New(schedule.SecondClosingTimeTuesday),
		SecondOpeningTimeWednesday: timestamppb.New(schedule.SecondOpeningTimeWednesday),
		SecondClosingTimeWednesday: timestamppb.New(schedule.SecondClosingTimeWednesday),
		SecondOpeningTimeThursday:  timestamppb.New(schedule.SecondOpeningTimeThursday),
		SecondClosingTimeThursday:  timestamppb.New(schedule.SecondClosingTimeThursday),
		SecondOpeningTimeFriday:    timestamppb.New(schedule.SecondOpeningTimeFriday),
		SecondClosingTimeFriday:    timestamppb.New(schedule.SecondClosingTimeFriday),
		SecondOpeningTimeSaturday:  timestamppb.New(schedule.SecondOpeningTimeSaturday),
		SecondClosingTimeSaturday:  timestamppb.New(schedule.SecondClosingTimeSaturday),
	}}, nil
}

func (v *businessService) GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest, md *utils.ClientMetadata) (*pb.Business, error) {
	var businessRes *entity.Business
	var businessCollectionRes *[]entity.BusinessCollection
	var businessErr, businessCollectionErr error
	var businessCollectionResponse []*pb.BusinessCollection
	err := v.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := v.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		businessId := uuid.MustParse(req.Id)
		businessRes, businessErr = v.dao.NewBusinessRepository().GetBusinessWithDistance(tx, &entity.Business{ID: &businessId}, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)})
		if businessErr != nil && businessErr.Error() == "record not found" {
			return errors.New("business not found")
		} else if businessErr != nil {
			return businessErr
		}
		businessCollectionRes, businessCollectionErr = v.dao.NewBusinessCollectionRepository().ListBusinessCollection(tx, &entity.BusinessCollection{BusinessId: &businessId}, nil)
		if businessCollectionErr != nil {
			return businessCollectionErr
		}
		businessCollectionResponse = make([]*pb.BusinessCollection, 0, len(*businessCollectionRes))
		for _, e := range *businessCollectionRes {
			businessCollectionResponse = append(businessCollectionResponse, &pb.BusinessCollection{
				Id:         e.ID.String(),
				Name:       e.Name,
				BusinessId: e.BusinessId.String(),
				Index:      e.Index,
				CreateTime: timestamppb.New(e.CreateTime),
				UpdateTime: timestamppb.New(e.UpdateTime),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var highQualityPhotoUrl, lowQualityPhotoUrl, thumbnailUrl string
	highQualityPhotoUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.HighQualityPhoto
	lowQualityPhotoUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.LowQualityPhoto
	thumbnailUrl = v.config.BusinessAvatarBulkName + "/" + businessRes.Thumbnail
	return &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, ToPickUp: businessRes.ToPickUp, DeliveryPriceCup: businessRes.DeliveryPriceCup, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.Coords()[1], Longitude: businessRes.Coordinates.Coords()[0]}, HighQualityPhotoUrl: highQualityPhotoUrl, LowQualityPhotoUrl: lowQualityPhotoUrl, ThumbnailUrl: thumbnailUrl, Distance: businessRes.Distance, BusinessCollections: businessCollectionResponse, BusinessCategory: businessRes.BusinessCategory}, nil
}
