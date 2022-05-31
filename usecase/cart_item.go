package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	gp "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CartItemService interface {
	ListCartItem(ctx context.Context, req *pb.ListCartItemRequest, md *utils.ClientMetadata) (*pb.ListCartItemResponse, error)
	AddCartItem(ctx context.Context, req *pb.AddCartItemRequest, md *utils.ClientMetadata) (*pb.AddCartItemResponse, error)
	IsEmptyCartItem(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.IsEmptyCartItemResponse, error)
	ReduceCartItem(ctx context.Context, req *pb.ReduceCartItemRequest, md *utils.ClientMetadata) (*pb.ReduceCartItemResponse, error)
	DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest, md *utils.ClientMetadata) (*gp.Empty, error)
	EmptyCartItem(request *dto.EmptyCartItemRequest) error
}

type cartItemService struct {
	dao    repository.DAO
	config *utils.Config
}

func NewCartItemService(dao repository.DAO, config *utils.Config) CartItemService {
	return &cartItemService{dao: dao, config: config}
}

func (i *cartItemService) EmptyCartItem(request *dto.EmptyCartItemRequest) error {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listCartItemsRes, listCartItemsErr := i.dao.NewCartItemRepository().ListCartItemAll(tx, &models.CartItem{UserId: authorizationTokenRes.UserId})
		if listCartItemsErr != nil {
			return listCartItemsErr
		}
		itemFks := make([]uuid.UUID, 0, len(*listCartItemsRes))
		for _, item := range *listCartItemsRes {
			itemFks = append(itemFks, *item.ItemId)
		}
		itemsRes, itemsErr := i.dao.NewItemQuery().ListItemInIds(tx, itemFks)
		if itemsErr != nil {
			return itemsErr
		}
		for _, item := range *listCartItemsRes {
			var index = -1
			for i, n := range *itemsRes {
				if n.ID == item.ItemId {
					index = i
				}
			}
			(*itemsRes)[index].Availability += int64(item.Quantity)
		}
		for _, item := range *itemsRes {
			_, updateItemsErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &item)
			if updateItemsErr != nil {
				return updateItemsErr
			}
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{UserId: authorizationTokenRes.UserId}, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *cartItemService) IsEmptyCartItem(ctx context.Context, req *gp.Empty, md *utils.ClientMetadata) (*pb.IsEmptyCartItemResponse, error) {
	var cartItemQuantityRes *bool
	var cartItemQuantityErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		cartItemQuantityRes, cartItemQuantityErr = i.dao.NewCartItemRepository().CartItemIsEmpty(tx, &models.CartItem{UserId: authorizationTokenRes.UserId})
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
	var items *[]models.CartItem
	var res pb.ListCartItemResponse
	var itemsErr error
	var nextPage time.Time
	if req.NextPage == nil {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		items, itemsErr = i.dao.NewCartItemRepository().ListCartItem(tx, &models.CartItem{UserId: authorizationTokenRes.UserId}, &nextPage)
		if itemsErr != nil {
			return itemsErr
		} else if len(*items) > 10 {
			*items = (*items)[:len(*items)-1]
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		} else if len(*items) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.CartItem, 0, len(*items))
	for _, item := range *items {
		itemsResponse = append(itemsResponse, &pb.CartItem{
			Id:                   item.ID.String(),
			Name:                 item.Name,
			Price:                item.Price,
			ItemId:               item.ItemId.String(),
			AuthorizationTokenId: item.AuthorizationTokenId.String(),
			Quantity:             item.Quantity,
			Thumbnail:            item.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + item.Thumbnail,
			ThumbnailBlurHash:    item.ThumbnailBlurHash,
			CreateTime:           timestamppb.New(item.CreateTime),
			UpdateTime:           timestamppb.New(item.UpdateTime),
		})
	}
	res.CartItems = itemsResponse
	return &res, nil
}

func (i *cartItemService) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest, md *utils.ClientMetadata) (*pb.AddCartItemResponse, error) {
	var result *models.CartItem
	var resultErr error
	itemId := uuid.MustParse(req.ItemId)
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, req.ItemId, location)
		var itemAvailability int64
		if itemErr != nil {
			return itemErr
		}
		if (item.Availability - int64(req.Quantity)) < 0 {
			return errors.New("no_availability:availability:" + strconv.Itoa(int(item.Availability)))
		} else if item.Availability-int64(req.Quantity) == 0 {
			itemAvailability = -1
		} else {
			itemAvailability = item.Availability - int64(req.Quantity)
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: itemAvailability})
		if updateItemErr != nil {
			return updateItemErr
		}
		cartItemRes, err := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: &itemId, UserId: authorizationTokenRes.UserId}, &[]string{"id", "quantity"})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if cartItemRes != nil {
			result, err = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemId: &itemId, UserId: authorizationTokenRes.UserId}, &models.CartItem{Quantity: cartItemRes.Quantity + req.Quantity})
			if err != nil {
				return err
			}
		} else if cartItemRes == nil && err.Error() == "record not found" {
			result, resultErr = i.dao.NewCartItemRepository().CreateCartItem(tx, &models.CartItem{Name: item.Name, Price: item.Price, Quantity: req.Quantity, ItemId: item.ID, UserId: authorizationTokenRes.UserId, AuthorizationTokenId: authorizationTokenRes.ID, BusinessId: item.BusinessId, Thumbnail: item.Thumbnail, ThumbnailBlurHash: item.ThumbnailBlurHash})
			if resultErr != nil {
				return resultErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddCartItemResponse{
		CartItem: &pb.CartItem{
			Id:                   result.ID.String(),
			Name:                 result.Name,
			Price:                result.Price,
			ItemId:               result.ItemId.String(),
			AuthorizationTokenId: result.AuthorizationTokenId.String(),
			Quantity:             result.Quantity,
			CreateTime:           timestamppb.New(result.CreateTime),
			UpdateTime:           timestamppb.New(result.UpdateTime),
			Thumbnail:            result.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + result.Thumbnail,
			ThumbnailBlurHash:    result.ThumbnailBlurHash,
		},
	}, nil
}

func (i *cartItemService) ReduceCartItem(ctx context.Context, req *pb.ReduceCartItemRequest, md *utils.ClientMetadata) (*pb.ReduceCartItemResponse, error) {
	var result *models.CartItem
	var resultErr error
	itemId := uuid.MustParse(req.ItemId)
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, req.ItemId, location)
		if itemErr != nil && itemErr.Error() == "record not found" {
			return errors.New("item not found")
		} else if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		result, resultErr = i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: &itemId, UserId: authorizationTokenRes.UserId}, &[]string{"id", "quantity"})
		if resultErr != nil && resultErr.Error() != "record not found" {
			return resultErr
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + 1})
		if updateItemErr != nil {
			return updateItemErr
		}
		if (result.Quantity - 1) == 0 {
			_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: result.ID, UserId: authorizationTokenRes.UserId}, nil)
			if err != nil {
				return err
			}
			result = nil
		} else {
			result, resultErr = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemId: &itemId, UserId: authorizationTokenRes.UserId}, &models.CartItem{Quantity: result.Quantity - 1})
			if resultErr != nil {
				return resultErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if result != nil {
		return &pb.ReduceCartItemResponse{CartItem: &pb.CartItem{
			Id:                   result.ID.String(),
			Name:                 result.Name,
			Price:                result.Price,
			ItemId:               result.ItemId.String(),
			AuthorizationTokenId: result.AuthorizationTokenId.String(),
			Quantity:             result.Quantity,
			CreateTime:           timestamppb.New(result.CreateTime),
			UpdateTime:           timestamppb.New(result.UpdateTime),
			Thumbnail:            result.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + result.Thumbnail,
			ThumbnailBlurHash:    result.ThumbnailBlurHash,
		}}, nil
	} else {
		return &pb.ReduceCartItemResponse{CartItem: &pb.CartItem{}}, nil
	}
}

func (i *cartItemService) DeleteCartItem(ctx context.Context, req *pb.DeleteCartItemRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	id := uuid.MustParse(req.Id)
	location := ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "user_id"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ID: &id, UserId: authorizationTokenRes.UserId}, &[]string{"id", "quantity", "item_id"})
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, cartItemRes.ItemId.String(), location)
		if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + int64(cartItemRes.Quantity)})
		if updateItemErr != nil {
			return updateItemErr
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: cartItemRes.ID, UserId: authorizationTokenRes.UserId}, nil)
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
