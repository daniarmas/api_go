package repository

import (
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartnerApplicationRepository interface {
	ListPartnerApplication(tx *gorm.DB, where *models.PartnerApplication, cursor *time.Time, fields *[]string) (*[]models.PartnerApplication, error)
	CreatePartnerApplication(tx *gorm.DB, where *models.PartnerApplication) (*models.PartnerApplication, error)
	UpdatePartnerApplication(tx *gorm.DB, where *models.PartnerApplication, data *models.PartnerApplication) (*models.PartnerApplication, error)
	GetPartnerApplication(tx *gorm.DB, where *models.PartnerApplication, fields *[]string) (*models.PartnerApplication, error)
	DeletePartnerApplication(tx *gorm.DB, where *models.PartnerApplication, ids *[]uuid.UUID) (*[]models.PartnerApplication, error)
}

type partnerApplicationRepository struct{}

func (i *partnerApplicationRepository) ListPartnerApplication(tx *gorm.DB, where *models.PartnerApplication, cursor *time.Time, fields *[]string) (*[]models.PartnerApplication, error) {
	res, err := Datasource.NewPartnerApplicationDatasource().ListPartnerApplication(tx, where, fields, cursor)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *partnerApplicationRepository) CreatePartnerApplication(tx *gorm.DB, data *models.PartnerApplication) (*models.PartnerApplication, error) {
	res, err := Datasource.NewPartnerApplicationDatasource().CreatePartnerApplication(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *partnerApplicationRepository) UpdatePartnerApplication(tx *gorm.DB, where *models.PartnerApplication, data *models.PartnerApplication) (*models.PartnerApplication, error) {
	res, err := Datasource.NewPartnerApplicationDatasource().UpdatePartnerApplication(tx, where, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *partnerApplicationRepository) GetPartnerApplication(tx *gorm.DB, where *models.PartnerApplication, fields *[]string) (*models.PartnerApplication, error) {
	res, err := Datasource.NewPartnerApplicationDatasource().GetPartnerApplication(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *partnerApplicationRepository) DeletePartnerApplication(tx *gorm.DB, where *models.PartnerApplication, ids *[]uuid.UUID) (*[]models.PartnerApplication, error) {
	res, err := Datasource.NewPartnerApplicationDatasource().DeletePartnerApplication(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
