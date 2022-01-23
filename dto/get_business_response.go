package dto

import "github.com/daniarmas/api_go/datastruct"

type GetBusinessResponse struct {
	Business     *datastruct.Business
	ItemCategory *[]datastruct.BusinessItemCategory
}
