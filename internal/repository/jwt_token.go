package repository

import (
	"time"

	"github.com/daniarmas/api_go/internal/datasource"
)

type JwtTokenRepository interface {
	CreateJwtRefreshToken(tokenMetadata *datasource.JsonWebTokenMetadata) error
	CreateJwtAccessToken(tokenMetadata *datasource.JsonWebTokenMetadata, expirationTime *time.Time) error
	CreateJwtAuthorizationToken(tokenMetadata *datasource.JsonWebTokenMetadata) error
	ParseJwtRefreshToken(tokenMetadata *datasource.JsonWebTokenMetadata) error
	ParseJwtAuthorizationToken(tokenMetadata *datasource.JsonWebTokenMetadata) error
}

type jwtTokenRepository struct{}

func (v *jwtTokenRepository) CreateJwtAccessToken(tokenMetadata *datasource.JsonWebTokenMetadata, expirationTime *time.Time) error {
	err := Datasource.NewJwtTokenDatasource().CreateJwtAccessToken(tokenMetadata, expirationTime)
	if err != nil {
		return err
	}
	return nil
}

func (v *jwtTokenRepository) CreateJwtRefreshToken(tokenMetadata *datasource.JsonWebTokenMetadata) error {
	err := Datasource.NewJwtTokenDatasource().CreateJwtRefreshToken(tokenMetadata)
	if err != nil {
		return err
	}
	return nil
}

func (r *jwtTokenRepository) CreateJwtAuthorizationToken(tokenMetadata *datasource.JsonWebTokenMetadata) error {
	err := Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(tokenMetadata)
	if err != nil {
		return err
	}
	return nil
}

func (r *jwtTokenRepository) ParseJwtRefreshToken(tokenMetadata *datasource.JsonWebTokenMetadata) error {
	err := Datasource.NewJwtTokenDatasource().ParseJwtRefreshToken(tokenMetadata)
	if err != nil {
		return err
	}
	return nil
}

func (r *jwtTokenRepository) ParseJwtAuthorizationToken(tokenMetadata *datasource.JsonWebTokenMetadata) error {
	err := Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(tokenMetadata)
	if err != nil {
		return err
	}
	return nil
}
