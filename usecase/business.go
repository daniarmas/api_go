package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type BusinessService interface {
	Feed(ctx context.Context, req *pb.FeedRequest, meta *utils.ClientMetadata) (*pb.FeedResponse, error)
	GetBusiness(ctx context.Context, req *pb.GetBusinessRequest, meta *utils.ClientMetadata) (*pb.GetBusinessResponse, error)
	GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest, md *utils.ClientMetadata) (*pb.GetBusinessWithDistanceResponse, error)
	CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest, md *utils.ClientMetadata) (*pb.CreateBusinessResponse, error)
	UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest, md *utils.ClientMetadata) (*pb.Business, error)
}

type businessService struct {
	config *utils.Config
	dao    repository.DAO
	stDb   *sql.DB
}

func NewBusinessService(dao repository.DAO, config *utils.Config, stDb *sql.DB) BusinessService {
	return &businessService{dao: dao, config: config, stDb: stDb}
}

func (i *businessService) UpdateBusiness(ctx context.Context, req *pb.UpdateBusinessRequest, md *utils.ClientMetadata) (*pb.Business, error) {
	var businessRes *models.Business
	var businessErr error
	id := uuid.MustParse(req.Id)
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &models.BusinessUser{UserId: authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		businessIsOpenRes, businessIsOpenErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: &id}, "OrderTypePickUp")
		if businessIsOpenErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessIsOpenErr
		} else if businessIsOpenRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: &id}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil && businessIsOpenErr.Error() != "business closed" {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{BusinessId: &id}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		getBusinessRes, getBusinessErr := i.dao.NewBusinessQuery().GetBusiness(tx, &models.Business{ID: &id}, nil)
		if getBusinessErr != nil {
			return getBusinessErr
		}
		if req.HighQualityPhoto != "" || req.LowQualityPhoto != "" || req.Thumbnail != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.HighQualityPhoto)
			if hqErr != nil && hqErr.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.LowQualityPhoto)
			if lqErr != nil && lqErr.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.Thumbnail)
			if tnErr != nil && tnErr.Error() == "ObjectMissing" {
				return errors.New("ThumbnailObject missing")
			} else if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.HighQualityPhoto})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.LowQualityPhoto})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getBusinessRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.BusinessAvatarBulkName, Object: getBusinessRes.Thumbnail})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.BusinessAvatarBulkName, getBusinessRes.Thumbnail, minio.RemoveObjectOptions{})
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
		businessRes, businessErr = i.dao.NewBusinessQuery().UpdateBusiness(tx, &models.Business{
			Name:                  req.Name,
			Address:               req.Address,
			HighQualityPhoto:      datasource.Config.BusinessAvatarBulkName + "/" + req.HighQualityPhoto,
			LowQualityPhoto:       datasource.Config.BusinessAvatarBulkName + "/" + req.LowQualityPhoto,
			Thumbnail:             datasource.Config.BusinessAvatarBulkName + "/" + req.Thumbnail,
			BlurHash:              req.BlurHash,
			TimeMarginOrderMonth:  req.TimeMarginOrderMonth,
			TimeMarginOrderDay:    req.TimeMarginOrderDay,
			TimeMarginOrderHour:   req.TimeMarginOrderHour,
			TimeMarginOrderMinute: req.TimeMarginOrderMinute,
			DeliveryPrice:         req.DeliveryPrice,
			ToPickUp:              req.ToPickUp,
			HomeDelivery:          req.HomeDelivery,
			ProvinceId:            &provinceId,
			MunicipalityId:        &municipalityId,
		}, &models.Business{ID: &id})
		if businessErr != nil {
			return businessErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, DeliveryPrice: businessRes.DeliveryPrice, TimeMarginOrderMonth: businessRes.TimeMarginOrderMonth, TimeMarginOrderDay: businessRes.TimeMarginOrderDay, TimeMarginOrderHour: businessRes.TimeMarginOrderHour, TimeMarginOrderMinute: businessRes.TimeMarginOrderMinute, ToPickUp: businessRes.ToPickUp, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), CreateTime: timestamppb.New(businessRes.CreateTime), UpdateTime: timestamppb.New(businessRes.UpdateTime)}, nil
}

func (i *businessService) CreateBusiness(ctx context.Context, req *pb.CreateBusinessRequest, md *utils.ClientMetadata) (*pb.CreateBusinessResponse, error) {
	var businessRes *models.Business
	var businessErr error
	var response pb.CreateBusinessResponse
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		businessOwnerRes, businessOwnerErr := i.dao.NewBusinessUserRepository().GetBusinessUser(tx, &models.BusinessUser{UserId: authorizationTokenRes.UserId}, nil)
		if businessOwnerErr != nil {
			return businessOwnerErr
		}
		if !businessOwnerRes.IsBusinessOwner {
			return errors.New("permission denied")
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("HighQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("LowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.BusinessAvatarBulkName, req.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("ThumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		provinceId := uuid.MustParse(req.ProvinceId)
		municipalityId := uuid.MustParse(req.MunicipalityId)
		businessBrandId := uuid.MustParse(req.BusinessBrandId)
		businessRes, businessErr = i.dao.NewBusinessQuery().CreateBusiness(tx, &models.Business{
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
			DeliveryPrice:         req.DeliveryPrice,
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
		response.Business = &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, DeliveryPrice: businessRes.DeliveryPrice, TimeMarginOrderMonth: businessRes.TimeMarginOrderMonth, TimeMarginOrderDay: businessRes.TimeMarginOrderDay, TimeMarginOrderHour: businessRes.TimeMarginOrderHour, TimeMarginOrderMinute: businessRes.TimeMarginOrderMinute, ToPickUp: businessRes.ToPickUp, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), CreateTime: timestamppb.New(businessRes.CreateTime), UpdateTime: timestamppb.New(businessRes.UpdateTime), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.FlatCoords()[0], Longitude: businessRes.Coordinates.FlatCoords()[1]}}
		var unionBusinessAndMunicipalities = make([]*models.UnionBusinessAndMunicipality, 0, len(req.Municipalities))
		for _, item := range req.Municipalities {
			municipalityId := uuid.MustParse(item)
			unionBusinessAndMunicipalities = append(unionBusinessAndMunicipalities, &models.UnionBusinessAndMunicipality{
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
	var businessRes *[]models.Business
	var businessResAdd *[]models.Business
	var businessErr, businessErrAdd error
	var response pb.FeedResponse
	if req.SearchMunicipalityType == pb.SearchMunicipalityType_More {
		err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, req.NextPage, false, req.HomeDelivery, req.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if len(*businessRes) > 5 {
			*businessRes = (*businessRes)[:len(*businessRes)-1]
			response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
			response.SearchMunicipalityType = pb.SearchMunicipalityType_More
		} else if len(*businessRes) <= 5 && len(*businessRes) != 0 {
			length := 5 - len(*businessRes)
			err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
				businessResAdd, businessErrAdd = v.dao.NewBusinessQuery().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, int32(length), req.ProvinceId, req.MunicipalityId, 0, true, req.HomeDelivery, req.ToPickUp)
				if businessErrAdd != nil {
					return businessErrAdd
				}
				return nil
			})
			if err != nil {
				return nil, err
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
			err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
				businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, 0, true, req.HomeDelivery, req.ToPickUp)
				if businessErr != nil {
					return businessErr
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if len(*businessRes) > 5 {
				*businessRes = (*businessRes)[:len(*businessRes)-1]
				response.NextPage = int32((*businessRes)[len(*businessRes)-1].Cursor)
				response.SearchMunicipalityType = pb.SearchMunicipalityType_More
			}
		}
	} else {
		err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
			businessRes, businessErr = v.dao.NewBusinessQuery().Feed(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, 5, req.ProvinceId, req.MunicipalityId, req.NextPage, true, req.HomeDelivery, req.ToPickUp)
			if businessErr != nil {
				return businessErr
			}
			return nil
		})
		if err != nil {
			return nil, err
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
				DeliveryPrice:         e.DeliveryPrice,
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
	return &response, nil
}

func (v *businessService) GetBusiness(ctx context.Context, req *pb.GetBusinessRequest, meta *utils.ClientMetadata) (*pb.GetBusinessResponse, error) {
	var businessRes *models.Business
	var businessCollectionRes *[]models.BusinessCollection
	var businessErr, itemCategoryErr error
	var itemsCategoryResponse []*pb.ItemCategory
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		businessId := uuid.MustParse(req.Id)
		businessRes, businessErr = v.dao.NewBusinessQuery().GetBusiness(tx, &models.Business{ID: &businessId}, nil)
		if businessErr != nil && businessErr.Error() == "record not found" {
			return errors.New("business not found")
		} else if businessErr != nil {
			return businessErr
		}
		businessCollectionRes, itemCategoryErr = v.dao.NewBusinessCollectionQuery().ListBusinessCollection(tx, &models.BusinessCollection{BusinessId: &businessId}, nil)
		if itemCategoryErr != nil {
			return itemCategoryErr
		}
		itemsCategoryResponse = make([]*pb.ItemCategory, 0, len(*businessCollectionRes))
		for _, e := range *businessCollectionRes {
			itemsCategoryResponse = append(itemsCategoryResponse, &pb.ItemCategory{
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
	return &pb.GetBusinessResponse{Business: &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, ToPickUp: businessRes.ToPickUp, DeliveryPrice: businessRes.DeliveryPrice, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.Coords()[1], Longitude: businessRes.Coordinates.Coords()[0]}, HighQualityPhotoUrl: highQualityPhotoUrl, LowQualityPhotoUrl: lowQualityPhotoUrl, ThumbnailUrl: thumbnailUrl}, ItemCategory: itemsCategoryResponse}, nil
}

func (v *businessService) GetBusinessWithDistance(ctx context.Context, req *pb.GetBusinessWithDistanceRequest, md *utils.ClientMetadata) (*pb.GetBusinessWithDistanceResponse, error) {
	var businessRes *models.Business
	var businessCollectionRes *[]models.BusinessCollection
	var businessErr, itemCategoryErr error
	var itemsCategoryResponse []*pb.ItemCategory
	// Collecting analytics
	if *md.App == "App" {
		go func() {
			// type TransactionPriority struct {
			// 	Priority string
			// }
			// var priority TransactionPriority
			ctx := context.Background()
			// Get a Tx for making transaction requests.
			tx, err := v.stDb.BeginTx(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}
			// Defer a rollback in case anything fails.
			defer tx.Rollback()

			// Set transaction priority
			_, err = tx.ExecContext(ctx, "SET TRANSACTION PRIORITY LOW")
			if err != nil {
				log.Fatal(err)
			}
			time := time.Now()
			_, err = tx.Exec(`INSERT INTO "business_analytics" ("type", "business_id", "create_time", "update_time") VALUES ($1, $2, $3, $4)`, "BusinessAnalyticsTypeDetailView", req.Id, time, time)
			if err != nil {
				log.Fatal(err)
			}

			// err = tx.QueryRow("SHOW transaction_priority").Scan(&priority.Priority)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// fmt.Println("TRANSACTION PRIORITY: " + priority.Priority)

			// Commit the transaction.
			if err = tx.Commit(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		businessId := uuid.MustParse(req.Id)
		businessRes, businessErr = v.dao.NewBusinessQuery().GetBusinessWithDistance(tx, &models.Business{ID: &businessId, Coordinates: ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}})
		if businessErr != nil && businessErr.Error() == "record not found" {
			return errors.New("business not found")
		} else if businessErr != nil {
			return businessErr
		}
		businessCollectionRes, itemCategoryErr = v.dao.NewBusinessCollectionQuery().ListBusinessCollection(tx, &models.BusinessCollection{BusinessId: &businessId}, nil)
		if itemCategoryErr != nil {
			return itemCategoryErr
		}
		itemsCategoryResponse = make([]*pb.ItemCategory, 0, len(*businessCollectionRes))
		for _, e := range *businessCollectionRes {
			itemsCategoryResponse = append(itemsCategoryResponse, &pb.ItemCategory{
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
	return &pb.GetBusinessWithDistanceResponse{Business: &pb.Business{Id: businessRes.ID.String(), Name: businessRes.Name, Address: businessRes.Address, HighQualityPhoto: businessRes.HighQualityPhoto, LowQualityPhoto: businessRes.LowQualityPhoto, Thumbnail: businessRes.Thumbnail, BlurHash: businessRes.BlurHash, ToPickUp: businessRes.ToPickUp, DeliveryPrice: businessRes.DeliveryPrice, HomeDelivery: businessRes.HomeDelivery, ProvinceId: businessRes.ProvinceId.String(), MunicipalityId: businessRes.MunicipalityId.String(), BusinessBrandId: businessRes.BusinessBrandId.String(), Coordinates: &pb.Point{Latitude: businessRes.Coordinates.Coords()[1], Longitude: businessRes.Coordinates.Coords()[0]}, HighQualityPhotoUrl: highQualityPhotoUrl, LowQualityPhotoUrl: lowQualityPhotoUrl, ThumbnailUrl: thumbnailUrl, Distance: businessRes.Distance}, ItemCategory: itemsCategoryResponse, BusinessCategory: businessRes.BusinessCategory}, nil
}
