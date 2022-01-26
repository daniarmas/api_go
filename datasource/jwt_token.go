package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtTokenDatasource interface {
	CreateJwtRefreshToken(refreshTokenFk *string) (*string, error)
	CreateJwtAuthorizationToken(authorizationTokenFk *string) (*string, error)
	ParseJwtRefreshToken(tokenValue *string) (*string, error)
	ParseJwtAuthorizationToken(tokenValue *string) (*string, error)
}

type jwtTokenDatasource struct{}

func (v *jwtTokenDatasource) CreateJwtRefreshToken(refreshTokenFk *string) (*string, error) {
	hmacSecret := []byte(Config.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		Subject:   *refreshTokenFk,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (r *jwtTokenDatasource) CreateJwtAuthorizationToken(authorizationTokenFk *string) (*string, error) {
	hmacSecret := []byte(Config.JwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   *authorizationTokenFk,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (r *jwtTokenDatasource) ParseJwtRefreshToken(tokenValue *string) (*string, error) {
	hmacSecret := []byte(Config.JwtSecret)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(*tokenValue, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tokenErr := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(tokenErr)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := fmt.Sprintf("%s", claims["sub"])
		return &data, nil
	} else {
		return nil, err
	}
}

func (r *jwtTokenDatasource) ParseJwtAuthorizationToken(tokenValue *string) (*string, error) {
	hmacSecret := []byte(Config.JwtSecret)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(*tokenValue, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tokenErr := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(tokenErr)
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := fmt.Sprintf("%s", claims["sub"])
		return &data, nil
	} else {
		return nil, err
	}
}
