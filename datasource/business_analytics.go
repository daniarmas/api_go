package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessAnalyticsDatasource interface {
	CreateBusinessAnalytics(tx *sql.Tx, data *[]entity.BusinessAnalytics) (*[]entity.BusinessAnalytics, error)
	GetBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*entity.BusinessAnalytics, error)
	ListBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*[]entity.BusinessAnalytics, error)
}

type businessAnalyticsDatasource struct{}

func (i *businessAnalyticsDatasource) CreateBusinessAnalytics(tx *sql.Tx, data *[]entity.BusinessAnalytics) (*[]entity.BusinessAnalytics, error) {
	var values string
	var res []entity.BusinessAnalytics
	for index, item := range *data {
		if index != len(*data)-1 {
			values = values + fmt.Sprintf("('%s', '%s', '%v', '%v'), ", item.Type, item.BusinessId, item.CreateTime.Format(time.RFC3339), item.CreateTime.Format(time.RFC3339))
		} else {
			values = values + fmt.Sprintf("('%s', '%s', '%v', '%v') RETURNING *;", item.Type, item.BusinessId, item.CreateTime.Format(time.RFC3339), item.CreateTime.Format(time.RFC3339))
		}
	}
	query := fmt.Sprintf(`INSERT INTO "business_analytics" ("type", "business_id", "create_time", "update_time") VALUES %s`, values)
	result, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var businessAnalytics entity.BusinessAnalytics
		if err := result.Scan(&businessAnalytics.ID, &businessAnalytics.Type, &businessAnalytics.BusinessId, &businessAnalytics.CreateTime, &businessAnalytics.UpdateTime, &businessAnalytics.DeleteTime); err != nil {
			return nil, err
		}
		res = append(res, businessAnalytics)
	}
	return &res, nil
}

func (v *businessAnalyticsDatasource) GetBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*entity.BusinessAnalytics, error) {
	var res *entity.BusinessAnalytics
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

func (i *businessAnalyticsDatasource) ListBusinessAnalytics(tx *gorm.DB, where *entity.BusinessAnalytics, fields *[]string) (*[]entity.BusinessAnalytics, error) {
	var res []entity.BusinessAnalytics
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
