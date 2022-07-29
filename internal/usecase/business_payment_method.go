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

type BusinessPaymentMethodService interface {
	ListBusinessPaymentMethod(ctx context.Context, req *pb.ListBusinessPaymentMethodRequest, md *utils.ClientMetadata) (*pb.ListBusinessPaymentMethodResponse, error)
}

type BusinesspaymentMethodService struct {
	dao    repository.Repository
	config *config.Config
	rdb    *redis.Client
	sqldb  *sqldb.Sql
}

func NewBusinessPaymentMethodService(dao repository.Repository, config *config.Config, rdb *redis.Client, sqldb *sqldb.Sql) BusinessPaymentMethodService {
	return &BusinesspaymentMethodService{dao: dao, config: config, rdb: rdb, sqldb: sqldb}
}

func (i *BusinesspaymentMethodService) ListBusinessPaymentMethod(ctx context.Context, req *pb.ListBusinessPaymentMethodRequest, md *utils.ClientMetadata) (*pb.ListBusinessPaymentMethodResponse, error) {
	var res pb.ListBusinessPaymentMethodResponse
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
		result, err := i.dao.NewBusinessPaymentMethodRepository().ListBusinessPaymentMethod(tx, &entity.BusinessPaymentMethod{})
		if err != nil {
			return err
		}
		businessPaymentMethods := make([]*pb.BusinessPaymentMethod, 0, len(*result))
		for _, item := range *result {
			businessPaymentMethods = append(businessPaymentMethods, &pb.BusinessPaymentMethod{
				Id:              item.ID.String(),
				Name:            item.Name,
				Type:            *utils.ParsePaymentMethodType(&item.Type),
				Address:         item.Address,
				Enabled:         item.Enabled,
				BusinessId:      item.BusinessId.String(),
				PaymentMethodId: item.PaymentMethodId.String(),
				CreateTime:      timestamppb.New(item.CreateTime),
				UpdateTime:      timestamppb.New(item.UpdateTime),
			})
		}
		res.BusinessPaymentMethods = businessPaymentMethods
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
