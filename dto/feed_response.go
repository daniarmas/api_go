package dto

type FeedResponse struct {
	Businesses             *[]Business
	CartQuantity           int32
	SearchMunicipalityType string
	NextPage               int32
}
