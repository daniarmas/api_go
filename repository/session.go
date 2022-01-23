package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type SessionQuery interface {
	ListSession(tx *gorm.DB, session *models.Session, fields *[]string) (*[]models.Session, error)
}

type sessionQuery struct{}

func (v *sessionQuery) ListSession(tx *gorm.DB, session *models.Session, fields *[]string) (*[]models.Session, error) {
	var sessionResult *[]models.Session
	var result *gorm.DB
	if fields != nil {
		if session != nil {
			result = tx.Model(&models.AuthorizationToken{}).Where(session).Select("authorization_token.id, authorization_token.device_fk, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_id").Joins("left join device on device.id = authorization_token.device_fk").Find(&sessionResult)
		} else {
			result = tx.Model(&models.AuthorizationToken{}).Select("authorization_token.id, authorization_token.device_fk, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_id").Joins("left join device on device.id = authorization_token.device_fk").Find(&sessionResult)
		}
	} else {
		if session != nil {
			result = tx.Model(&models.AuthorizationToken{}).Where(session).Select("authorization_token.id, authorization_token.device_fk, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_id").Joins("left join device on device.id = authorization_token.device_fk").Find(&sessionResult)
		} else {
			result = tx.Model(&models.AuthorizationToken{}).Select("authorization_token.id, authorization_token.device_fk, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_id").Joins("left join device on device.id = authorization_token.device_fk").Find(&sessionResult)
		}
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return sessionResult, nil
		} else {
			return nil, result.Error
		}
	}
	return sessionResult, nil
}