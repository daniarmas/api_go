package service

import (
	"github.com/daniarmas/api_go/src/datastruct"
	"github.com/daniarmas/api_go/src/repository"
)

type AuthenticationService interface {
	CreateVerificationCode(verificationCode *datastruct.VerificationCode) error
	GetVerificationCode(code string, email string, verificationCodeType string, deviceId string) error
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) CreateVerificationCode(verificationCode *datastruct.VerificationCode) error {
	error := v.dao.NewVerificationCodeQuery().CreateVerificationCode(verificationCode)
	return error
}

func (v *authenticationService) GetVerificationCode(code string, email string, verificationCodeType string, deviceId string) error {
	error := v.dao.NewVerificationCodeQuery().GetVerificationCode(code, email, verificationCodeType, deviceId)
	return error
}
