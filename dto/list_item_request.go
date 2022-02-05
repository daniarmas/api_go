package dto

import "time"

type ListItemRequest struct {
	BusinessFk             string
	BusinessItemCategoryFk string
	NextPage               time.Time
}
