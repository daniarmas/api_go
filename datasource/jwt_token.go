package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JsonWebTokenMetadata struct {
	TokenId *uuid.UUID
	Token   *string
}

type JwtTokenDatasource interface {
	CreateJwtRefreshToken(tokenMetadata *JsonWebTokenMetadata) error
	CreateJwtAuthorizationToken(tokenMetadata *JsonWebTokenMetadata) error
	ParseJwtRefreshToken(tokenMetadata *JsonWebTokenMetadata) error
	ParseJwtAuthorizationToken(tokenMetadata *JsonWebTokenMetadata) error
}

type jwtTokenDatasource struct{}

func (v *jwtTokenDatasource) CreateJwtRefreshToken(tokenMetadata *JsonWebTokenMetadata) error {
	hmacSecret := []byte(Config.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		Subject:   tokenMetadata.TokenId.String(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return err
	}
	*tokenMetadata.Token = tokenString
	return nil
}

func (r *jwtTokenDatasource) CreateJwtAuthorizationToken(tokenMetadata *JsonWebTokenMetadata) error {
	hmacSecret := []byte(Config.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   tokenMetadata.TokenId.String(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return err
	}
	*tokenMetadata.Token = tokenString
	return nil
}

func (r *jwtTokenDatasource) ParseJwtRefreshToken(tokenMetadata *JsonWebTokenMetadata) error {
	hmacSecret := []byte(Config.JwtSecret)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(*tokenMetadata.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tokenErr := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(tokenErr)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := fmt.Sprintf("%s", claims["sub"])
		*tokenMetadata.TokenId = uuid.MustParse(data)
		return nil
	} else {
		return err
	}
}

func (r *jwtTokenDatasource) ParseJwtAuthorizationToken(tokenMetadata *JsonWebTokenMetadata) error {
	hmacSecret := []byte(Config.JwtSecret)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(*tokenMetadata.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tokenErr := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(tokenErr)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := fmt.Sprintf("%s", claims["sub"])
		*tokenMetadata.TokenId = uuid.MustParse(data)
		return nil
	} else {
		return err
	}
}
