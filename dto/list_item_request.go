package dto

type ListItemRequest struct {
	BusinessFk             string
	BusinessItemCategoryFk string
	NextPage               int32
}
