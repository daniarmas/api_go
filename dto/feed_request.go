package dto

type FeedRequest struct {
	ProvinceFk             string
	MunicipalityFk         string
	NextPage               int32
	Location               Point
	SearchMunicipalityType string
}

type Point struct {
	latitude  float32
	longitude float32
}
