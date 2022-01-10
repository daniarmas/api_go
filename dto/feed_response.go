package dto

type FeedResponse struct {
  Businesses *[]Business
  IsCache bool
  SearchMunicipalityType string 
  NextPage int32 
}