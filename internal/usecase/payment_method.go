package usecase

import (
	"context"
	"errors"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PaymentMethodService interface {
	UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest, md *utils.ClientMetadata) (*pb.PaymentMethod, error)
	CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, md *utils.ClientMetadata) (*pb.PaymentMethod, error)
	ListPaymentMethod(ctx context.Context, req *pb.ListPaymentMethodRequest, md *utils.ClientMetadata) (*pb.ListPaymentMethodResponse, error)
	DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type paymentMethodService struct {
	dao    repository.Repository
	config *config.Config
	rdb    *redis.Client
	sqldb  *sqldb.Sql
}

func NewPaymentMethodService(dao repository.Repository, config *config.Config, rdb *redis.Client, sqldb *sqldb.Sql) PaymentMethodService {
	return &paymentMethodService{dao: dao, config: config, rdb: rdb, sqldb: sqldb}
}

func (i *paymentMethodService) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
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
		_, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		// _, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_payment_method"}, &[]string{"id"})
		// if permissionErr != nil && permissionErr.Error() == "record not found" {
		// 	return errors.New("permission denied")
		// }
		id := uuid.MustParse(req.Id)
		_, err = i.dao.NewPaymentMethodRepository().DeletePaymentMethod(tx, &entity.PaymentMethod{ID: &id}, nil)
		if err != nil && err.Error() == "record not found" {
			return errors.New("payment method not found")
		} else if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &gp.Empty{}, nil
}

func (i *paymentMethodService) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest, md *utils.ClientMetadata) (*pb.PaymentMethod, error) {
	var res pb.PaymentMethod
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_payment_method"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		id := uuid.MustParse(req.Id)
		result, err := i.dao.NewPaymentMethodRepository().UpdatePaymentMethod(tx, &entity.PaymentMethod{ID: &id}, &entity.PaymentMethod{Name: req.PaymentMethod.Name, Enabled: req.PaymentMethod.Enabled, Address: req.PaymentMethod.Address, Type: req.PaymentMethod.Type.String()})
		if err != nil && err.Error() == "record not found" {
			return errors.New("payment method not found")
		} else if err != nil {
			return err
		}
		res = pb.PaymentMethod{
			Id:         result.ID.String(),
			Name:       result.Name,
			Type:       *utils.ParsePaymentMethodType(&result.Type),
			Address:    result.Address,
			Enabled:    result.Enabled,
			CreateTime: timestamppb.New(result.CreateTime),
			UpdateTime: timestamppb.New(result.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *paymentMethodService) ListPaymentMethod(ctx context.Context, req *pb.ListPaymentMethodRequest, md *utils.ClientMetadata) (*pb.ListPaymentMethodResponse, error) {
	var res pb.ListPaymentMethodResponse
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
		_, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		// _, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "list_payment_method"}, &[]string{"id"})
		// if permissionErr != nil && permissionErr.Error() == "record not found" {
		// 	return errors.New("permission denied")
		// }
		result, err := i.dao.NewPaymentMethodRepository().ListPaymentMethod(tx, &entity.PaymentMethod{})
		if err != nil {
			return err
		}
		paymentMethods := make([]*pb.PaymentMethod, 0, len(*result))
		for _, item := range *result {
			paymentMethods = append(paymentMethods, &pb.PaymentMethod{
				Id:         item.ID.String(),
				Name:       item.Name,
				Type:       *utils.ParsePaymentMethodType(&item.Type),
				Address:    item.Address,
				Enabled:    item.Enabled,
				CreateTime: timestamppb.New(item.CreateTime),
				UpdateTime: timestamppb.New(item.UpdateTime),
			})
		}
		res.PaymentMethods = paymentMethods
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *paymentMethodService) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, md *utils.ClientMetadata) (*pb.PaymentMethod, error) {
	var res pb.PaymentMethod
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
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("authorization token not found")
		} else if err != nil && err.Error() != "record not found" {
			return err
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "create_payment_method"}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		result, err := i.dao.NewPaymentMethodRepository().CreatePaymentMethod(tx, &entity.PaymentMethod{Name: req.PaymentMethod.Name, Enabled: req.PaymentMethod.Enabled, Address: req.PaymentMethod.Address, Type: req.PaymentMethod.Type.String()})
		if err != nil {
			return err
		}
		res = pb.PaymentMethod{
			Id:         result.ID.String(),
			Name:       result.Name,
			Type:       *utils.ParsePaymentMethodType(&result.Type),
			Address:    result.Address,
			Enabled:    result.Enabled,
			CreateTime: timestamppb.New(result.CreateTime),
			UpdateTime: timestamppb.New(result.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
