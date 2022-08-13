package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CartItemService interface {
	ListCartItem(ctx context.Context, req *pb.ListCartItemRequest, md *utils.ClientMetadata) (*pb.ListCartItemResponse, error)
	AddCartItem(ctx context.Context, req *pb.AddCartItemRequest, md *utils.ClientMetadata) (*pb.CartItem, error)
	EmptyAndAddCartItem(ctx context.Context, req *pb.EmptyAndAddCartItemRequest, md *utils.ClientMetadata) (*pb.CartItem, error)
	IsEmptyCartItem(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.IsEmptyCartItemResponse, error)
	DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	EmptyCartItem(ctx context.Context, md *utils.ClientMetadata) (*gp.Empty, error)
}

type cartItemService struct {
	dao    repository.Repository
	config *config.Config
	sqldb  *sqldb.Sql
}

func NewCartItemService(dao repository.Repository, config *config.Config, sqldb *sqldb.Sql) CartItemService {
	return &cartItemService{dao: dao, config: config, sqldb: sqldb}
}

func (i *cartItemService) EmptyAndAddCartItem(ctx context.Context, req *pb.EmptyAndAddCartItemRequest, md *utils.ClientMetadata) (*pb.CartItem, error) {
	var res pb.CartItem
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
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		listCartItemsRes, err := i.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
		if err != nil {
			return err
		}
		itemFks := make([]uuid.UUID, 0, len(*listCartItemsRes))
		for _, item := range *listCartItemsRes {
			itemFks = append(itemFks, *item.ItemId)
		}
		itemsRes, err := i.dao.NewItemRepository().ListItemInIds(ctx, tx, itemFks)
		if err != nil {
			return err
		}
		for _, item := range *listCartItemsRes {
			var index = -1
			for i, n := range *itemsRes {
				if *n.ID == *item.ItemId {
					index = i
				}
			}
			(*itemsRes)[index].Availability += int64(item.Quantity)
		}
		for _, item := range *itemsRes {
			_, err := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &item)
			if err != nil {
				return err
			}
		}
		_, err = i.dao.NewCartItemRepository().DeleteCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		// Add CartItem
		itemId := uuid.MustParse(req.ItemId)
		item, err := i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: &itemId})
		var itemAvailability int64
		if err != nil && err.Error() == "record not found" {
			return errors.New("item not found")
		} else if err != nil {
			return err
		}
		var result *entity.CartItem
		cartItemRes, err := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{ItemId: &itemId})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if cartItemRes != nil {
			// Restoring the itemAvailability
			item.Availability = item.Availability + int64(cartItemRes.Quantity)
			if (item.Availability - int64(req.Quantity)) < 0 {
				return errors.New("no_availability:availability:" + strconv.Itoa(int(item.Availability)))
			} else if item.Availability-int64(req.Quantity) == 0 {
				itemAvailability = -1
			} else {
				itemAvailability = item.Availability - int64(req.Quantity)
			}
			_, err := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &entity.Item{Availability: itemAvailability})
			if err != nil {
				return err
			}
			result, err = i.dao.NewCartItemRepository().UpdateCartItem(tx, &entity.CartItem{ItemId: &itemId}, &entity.CartItem{Quantity: req.Quantity})
			if err != nil {
				return err
			}
		} else if cartItemRes == nil && err.Error() == "record not found" {
			result, err = i.dao.NewCartItemRepository().CreateCartItem(tx, &entity.CartItem{Name: item.Name, PriceCup: item.PriceCup, Quantity: req.Quantity, ItemId: item.ID, UserId: authorizationTokenRes.UserId, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: item.BusinessId, Thumbnail: item.Thumbnail, BlurHash: item.BlurHash})
			if err != nil {
				return err
			}
			_, err = i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &entity.Item{Availability: item.Availability - int64(req.Quantity)})
			if err != nil {
				return err
			}
		}
		res = pb.CartItem{
			Id:                   result.ID.String(),
			Name:                 result.Name,
			PriceCup:             result.PriceCup,
			ItemId:               result.ItemId.String(),
			BusinessId:           result.BusinessId.String(),
			AuthorizationTokenId: result.AuthorizationTokenId.String(),
			Quantity:             result.Quantity,
			CreateTime:           timestamppb.New(result.CreateTime),
			UpdateTime:           timestamppb.New(result.UpdateTime),
			Thumbnail:            result.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + result.Thumbnail,
			BlurHash:             result.BlurHash,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *cartItemService) EmptyCartItem(ctx context.Context, md *utils.ClientMetadata) (*gp.Empty, error) {
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
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		listCartItemsRes, err := i.dao.NewCartItemRepository().ListCartItemAll(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
		if err != nil {
			return err
		}
		itemFks := make([]uuid.UUID, 0, len(*listCartItemsRes))
		for _, item := range *listCartItemsRes {
			itemFks = append(itemFks, *item.ItemId)
		}
		itemsRes, err := i.dao.NewItemRepository().ListItemInIds(ctx, tx, itemFks)
		if err != nil {
			return err
		}
		for _, item := range *listCartItemsRes {
			var index = -1
			for i, n := range *itemsRes {
				if *n.ID == *item.ItemId {
					index = i
				}
			}
			(*itemsRes)[index].Availability += int64(item.Quantity)
		}
		for _, item := range *itemsRes {
			_, err := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &item)
			if err != nil {
				return err
			}
		}
		_, err = i.dao.NewCartItemRepository().DeleteCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, nil)
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

func (i *cartItemService) IsEmptyCartItem(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.IsEmptyCartItemResponse, error) {
	var cartItemQuantityRes *bool
	var cartItemQuantityErr error
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		cartItemQuantityRes, cartItemQuantityErr = i.dao.NewCartItemRepository().CartItemIsEmpty(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
		if cartItemQuantityErr != nil {
			return cartItemQuantityErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.IsEmptyCartItemResponse{IsEmpty: *cartItemQuantityRes}, nil
}

func (i *cartItemService) ListCartItem(ctx context.Context, req *pb.ListCartItemRequest, md *utils.ClientMetadata) (*pb.ListCartItemResponse, error) {
	var res pb.ListCartItemResponse
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
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		var nextPage time.Time
		if req.NextPage == nil {
			nextPage = time.Now()
		} else {
			nextPage = req.NextPage.AsTime()
		}
		items, err := i.dao.NewCartItemRepository().ListCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId}, &nextPage)
		if err != nil {
			return err
		} else if len(*items) > 10 {
			*items = (*items)[:len(*items)-1]
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		} else if len(*items) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		}
		itemsResponse := make([]*pb.CartItem, 0, len(*items))
		for _, item := range *items {
			itemsResponse = append(itemsResponse, &pb.CartItem{
				Id:                   item.ID.String(),
				Name:                 item.Name,
				PriceCup:             item.PriceCup,
				ItemId:               item.ItemId.String(),
				BusinessId:           item.BusinessId.String(),
				AuthorizationTokenId: item.AuthorizationTokenId.String(),
				Quantity:             item.Quantity,
				Thumbnail:            item.Thumbnail,
				ThumbnailUrl:         i.config.ItemsBulkName + "/" + item.Thumbnail,
				BlurHash:             item.BlurHash,
				CreateTime:           timestamppb.New(item.CreateTime),
				UpdateTime:           timestamppb.New(item.UpdateTime),
			})
		}
		res.CartItems = itemsResponse
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *cartItemService) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest, md *utils.ClientMetadata) (*pb.CartItem, error) {
	var res pb.CartItem
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
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		itemId := uuid.MustParse(req.ItemId)
		item, err := i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: &itemId})
		var itemAvailability int64
		if err != nil && err.Error() == "record not found" {
			return errors.New("item not found")
		} else if err != nil {
			return err
		}
		var result *entity.CartItem
		cartItemRes, err := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{ItemId: &itemId})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if cartItemRes != nil {
			// Restoring the itemAvailability
			item.Availability = item.Availability + int64(cartItemRes.Quantity)
			if (item.Availability - int64(req.Quantity)) < 0 {
				return errors.New("no_availability:availability:" + strconv.Itoa(int(item.Availability)))
			} else if item.Availability-int64(req.Quantity) == 0 {
				itemAvailability = -1
			} else {
				itemAvailability = item.Availability - int64(req.Quantity)
			}
			_, err = i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &entity.Item{Availability: itemAvailability})
			if err != nil {
				return err
			}
			result, err = i.dao.NewCartItemRepository().UpdateCartItem(tx, &entity.CartItem{ItemId: &itemId}, &entity.CartItem{Quantity: req.Quantity})
			if err != nil {
				return err
			}
		} else if cartItemRes == nil && err.Error() == "record not found" {
			cartItemExists, err := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{UserId: authorizationTokenRes.UserId})
			if err != nil && err.Error() != "record not found" {
				return err
			} else if cartItemExists != nil && *cartItemExists.BusinessId != *item.BusinessId {
				return errors.New("the items in the cart can only be from one business")
			}
			result, err = i.dao.NewCartItemRepository().CreateCartItem(tx, &entity.CartItem{Name: item.Name, PriceCup: item.PriceCup, Quantity: req.Quantity, ItemId: item.ID, UserId: authorizationTokenRes.UserId, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: item.BusinessId, Thumbnail: item.Thumbnail, BlurHash: item.BlurHash})
			if err != nil {
				return err
			}
			_, err = i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &entity.Item{Availability: item.Availability - int64(req.Quantity)})
			if err != nil {
				return err
			}
		}
		res = pb.CartItem{
			Id:                   result.ID.String(),
			Name:                 result.Name,
			PriceCup:             result.PriceCup,
			ItemId:               result.ItemId.String(),
			BusinessId:           result.BusinessId.String(),
			AuthorizationTokenId: result.AuthorizationTokenId.String(),
			Quantity:             result.Quantity,
			CreateTime:           timestamppb.New(result.CreateTime),
			UpdateTime:           timestamppb.New(result.UpdateTime),
			Thumbnail:            result.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + result.Thumbnail,
			BlurHash:             result.BlurHash,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *cartItemService) DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		var whereCartItem entity.CartItem
		if req.ItemId != "" {
			value := uuid.MustParse(req.ItemId)
			whereCartItem.ItemId = &value
		}
		if req.Id != "" {
			value := uuid.MustParse(req.Id)
			whereCartItem.ID = &value
		}
		whereCartItem.UserId = authorizationTokenRes.UserId
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &whereCartItem)
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		item, itemErr := i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: cartItemRes.ItemId})
		if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		_, updateItemErr := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: item.ID}, &entity.Item{Availability: item.Availability + int64(cartItemRes.Quantity)})
		if updateItemErr != nil {
			return updateItemErr
		}
		_, err = i.dao.NewCartItemRepository().DeleteCartItem(tx, &whereCartItem, nil)
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
