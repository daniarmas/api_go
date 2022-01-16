package dto

import "github.com/daniarmas/api_go/datastruct"

type SearchItemResponse struct {
	Items                  *[]datastruct.Item
	SearchMunicipalityType string
	NextPage               int32
}
