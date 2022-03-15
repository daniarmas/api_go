package repository

import (
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type BusinessQuery interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error)
	CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error)
	GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error)
	GetBusinessWithLocation(tx *gorm.DB, where *models.Business) (*models.Business, error)
	GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error)
	UpdateBusiness(tx *gorm.DB, data *models.Business, where *models.Business) (*models.Business, error)
	UpdateBusinessCoordinate(tx *gorm.DB, data *models.Business, where *models.Business) error
}

type businessQuery struct{}

func (b *businessQuery) CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error) {
	res, err := Datasource.NewBusinessDatasource().CreateBusiness(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *businessQuery) UpdateBusiness(tx *gorm.DB, data *models.Business, where *models.Business) (*models.Business, error) {
	res, err := Datasource.NewBusinessDatasource().UpdateBusiness(tx, data, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *businessQuery) UpdateBusinessCoordinate(tx *gorm.DB, data *models.Business, where *models.Business) error {
	err := Datasource.NewBusinessDatasource().UpdateBusinessCoordinate(tx, data, where)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessQuery) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceFk string, municipalityFk string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().Feed(tx, coordinates, limit, provinceFk, municipalityFk, cursor, municipalityNotEqual, homeDelivery, toPickUp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessQuery) GetBusinessWithLocation(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusinessWithLocation(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessQuery) GetBusiness(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusiness(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessQuery) GetBusinessProvinceAndMunicipality(tx *gorm.DB, businessFk uuid.UUID) (*dto.GetBusinessProvinceAndMunicipality, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusinessProvinceAndMunicipality(tx, businessFk)
	if err != nil {
		return nil, err
	}
	return result, nil
}
