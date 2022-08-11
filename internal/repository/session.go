package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type SessionRepository interface {
	ListSession(tx *gorm.DB, session *entity.Session) (*[]entity.Session, error)
}

type sessionRepository struct{}

func (v *sessionRepository) ListSession(tx *gorm.DB, session *entity.Session) (*[]entity.Session, error) {
	var sessionResult *[]entity.Session
	var result *gorm.DB
	if session != nil {
		result = tx.Model(&entity.AuthorizationToken{}).Where(session).Select("authorization_token.id, authorization_token.device_id, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_identifier").Joins("left join device on device.id = authorization_token.device_id").Find(&sessionResult)
	} else {
		result = tx.Model(&entity.AuthorizationToken{}).Select("authorization_token.id, authorization_token.device_id, authorization_token.app, authorization_token.app_version, device.platform, device.system_version, device.model, device.device_identifier").Joins("left join device on device.id = authorization_token.device_id").Find(&sessionResult)
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
