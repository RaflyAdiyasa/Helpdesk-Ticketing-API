package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTservice interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct {
	secret string
	expiry time.Duration
}

func NewJWTservice(secret string, expiry time.Duration) JWTservice {
	return &jwtService{
		secret: secret,
		expiry: expiry,
	}
}

func (s *jwtService) GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.expiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected sighning method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
