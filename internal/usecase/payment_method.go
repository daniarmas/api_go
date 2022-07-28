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
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PaymentMethodService interface {
	CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, md *utils.ClientMetadata) (*pb.PaymentMethod, error)
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
