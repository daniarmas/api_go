package usecase

import (
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error)
	// CreateOrder(request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
}

type orderService struct {
	dao repository.DAO
}

func NewOrderService(dao repository.DAO) OrderService {
	return &orderService{dao: dao}
}

func (i *orderService) ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error) {
	var listOrderResponse dto.ListOrderResponse
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		ordersRes, ordersErr := i.dao.NewOrderRepository().ListOrder(tx, &models.Order{CreateTime: request.NextPage, UserFk: authorizationTokenRes.UserFk})
		if ordersErr != nil {
			return ordersErr
		}
		if len(*ordersRes) > 10 {
			*ordersRes = (*ordersRes)[:len(*ordersRes)-1]
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else if len(*ordersRes) > 0 {
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else {
			listOrderResponse.NextPage = request.NextPage
		}
		listOrderResponse.Orders = ordersRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &listOrderResponse, nil
}
