package usecase

import (
	"context"
	"database/sql"

	"github.com/daniarmas/api_go/internal/entity"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsService interface {
	CollectAnalytics(ctx context.Context, req *pb.CollectAnalyticsRequest, md *utils.ClientMetadata) (*gp.Empty, error)
}

type analyticsService struct {
	dao  repository.Repository
	stDb *sql.DB
}

func NewAnalyticsService(dao repository.Repository, stDb *sql.DB) AnalyticsService {
	return &analyticsService{dao: dao, stDb: stDb}
}

func (i *analyticsService) CollectAnalytics(ctx context.Context, req *pb.CollectAnalyticsRequest, md *utils.ClientMetadata) (*gp.Empty, error) {
	// Collecting analytics
	if *md.App == "App" {
		go func() {
			ctx := context.Background()
			// Get a Tx for making transaction requests.
			tx, err := i.stDb.BeginTx(ctx, nil)
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

			var businessAnalytics []entity.BusinessAnalytics
			for _, i := range req.BusinessAnalytics {
				businessId := uuid.MustParse(i.BusinessId)
				businessAnalytics = append(businessAnalytics, entity.BusinessAnalytics{
					Type:       i.Type.String(),
					BusinessId: &businessId,
					CreateTime: i.CreateTime.AsTime(),
					UpdateTime: i.CreateTime.AsTime(),
				})
			}
			var itemAnalytics []entity.ItemAnalytics
			for _, i := range req.ItemAnalytics {
				itemId := uuid.MustParse(i.ItemId)
				itemAnalytics = append(itemAnalytics, entity.ItemAnalytics{
					Type:       i.Type.String(),
					ItemId:     &itemId,
					CreateTime: i.CreateTime.AsTime(),
					UpdateTime: i.CreateTime.AsTime(),
				})
			}
			_, err = i.dao.NewBusinessAnalyticsRepository().CreateBusinessAnalytics(tx, &businessAnalytics)
			if err != nil {
				log.Fatal(err)
			}
			_, err = i.dao.NewItemAnalyticsRepository().CreateItemAnalytics(tx, &itemAnalytics)
			if err != nil {
				log.Fatal(err)
			}

			// Commit the transaction.
			if err = tx.Commit(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	return &gp.Empty{}, nil
}
