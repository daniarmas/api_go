package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type ItemAnalyticsDatasource interface {
	CreateItemAnalytics(tx *sql.Tx, data *[]models.ItemAnalytics) (*[]models.ItemAnalytics, error)
	GetItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*models.ItemAnalytics, error)
	ListItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*[]models.ItemAnalytics, error)
}

type itemAnalyticsDatasource struct{}

func (i *itemAnalyticsDatasource) CreateItemAnalytics(tx *sql.Tx, data *[]models.ItemAnalytics) (*[]models.ItemAnalytics, error) {
	var values string
	var res []models.ItemAnalytics
	for index, item := range *data {
		if index != len(*data)-1 {
			values = values + fmt.Sprintf("('%s', '%s', '%v', '%v'), ", item.Type, item.ItemId, item.CreateTime.Format(time.RFC3339), item.CreateTime.Format(time.RFC3339))
		} else {
			values = values + fmt.Sprintf("('%s', '%s', '%v', '%v') RETURNING *;", item.Type, item.ItemId, item.CreateTime.Format(time.RFC3339), item.CreateTime.Format(time.RFC3339))
		}
	}
	query := fmt.Sprintf(`INSERT INTO "item_analytics" ("type", "item_id", "create_time", "update_time") VALUES %s`, values)
	result, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var itemAnalytics models.ItemAnalytics
		if err := result.Scan(&itemAnalytics.ID, &itemAnalytics.Type, &itemAnalytics.ItemId, &itemAnalytics.CreateTime, &itemAnalytics.UpdateTime, &itemAnalytics.DeleteTime); err != nil {
			return nil, err
		}
		res = append(res, itemAnalytics)
	}
	return &res, nil
}

func (v *itemAnalyticsDatasource) GetItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*models.ItemAnalytics, error) {
	var res *models.ItemAnalytics
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (i *itemAnalyticsDatasource) ListItemAnalytics(tx *gorm.DB, where *models.ItemAnalytics, fields *[]string) (*[]models.ItemAnalytics, error) {
	var res []models.ItemAnalytics
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
