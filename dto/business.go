package dto

type Business struct {
  ID string
  Name string
  Description string
  Address string
  Phone string
  Email string
  highQualityPhoto string
  highQualityPhotoBlurHash string
  LowQualityPhoto string
  LowQualityPhotoBlurHash string
  Thumbnail string
  ThumbnailBlurHash string
  IsOpen bool
  DeliveryPrice float32
  Polygon []Polygon
  Coordinates Point
  LeadDayTime int32
  LeadHoursTime int32
  LeadMinutesTime int32
  ToPickUp bool
  HomeDelivery bool
  BusinessBrandFk string
  ProvinceFk string
  MunicipalityFk string
  Distance float32
  Status string 
  cursor int32
}