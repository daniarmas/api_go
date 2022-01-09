package service

import (
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/repository"
)

type BusinessService interface {
	Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error)
}

type businessService struct {
	dao repository.DAO
}

func NewBusinessService(dao repository.DAO) BusinessService {
	return &businessService{dao: dao}
}

func (b *businessService) Feed(feedRequest *dto.FeedRequest) (*dto.FeedResponse, error) {
	return nil, nil
}
