package service

import (
	"github.com/daniarmas/api_go/datastruct"
	"github.com/daniarmas/api_go/dto"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"gorm.io/gorm"
)

type ItemService interface {
	GetItem(id string) (*datastruct.Item, error)
	ListItem() ([]datastruct.Item, error)
	SearchItem(name string, provinceFk string, municipalityFk string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error)
	// CreateItem(answer datastruct.Item) (*int64, error)
	// UpdateItem(answer datastruct.Item) (*datastruct.Item, error)
	// DeleteItem(id int64) error
}

type itemService struct {
	dao repository.DAO
}

func NewItemService(dao repository.DAO) ItemService {
	return &itemService{dao: dao}
}

func (i *itemService) ListItem() ([]datastruct.Item, error) {
	var items []datastruct.Item
	var itemsErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		items, itemsErr = i.dao.NewItemQuery().ListItem(tx)
		if itemsErr != nil {
			return itemsErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *itemService) GetItem(id string) (*datastruct.Item, error) {
	item, err := i.dao.NewItemQuery().GetItem(id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (i *itemService) SearchItem(name string, provinceFk string, municipalityFk string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error) {
	var response *[]datastruct.Item
	var searchItemResponse dto.SearchItemResponse
	var responseErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		if searchMunicipalityType == "More" {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, false, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More.String()
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, int64(length))
				if responseErr != nil {
					return responseErr
				}
				if len(*responseAdd) > length {
					*responseAdd = (*responseAdd)[:len(*responseAdd)-1]
				}
				*response = append(*response, *responseAdd...)
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
			} else if len(*response) == 0 {
				response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, 10)
				if responseErr != nil {
					return responseErr
				}
				if len(*response) > 10 {
					*response = (*response)[:len(*response)-1]
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				} else if len(*response) <= 10 && len(*response) != 0 {
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				}
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
			}
		} else {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			} else if len(*response) <= 10 && len(*response) != 0 {
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			}
			searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	searchItemResponse.Items = response
	return &searchItemResponse, nil
}
