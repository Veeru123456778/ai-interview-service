package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret []byte
}

func NewService(secret string) *Service {
	return &Service{
		jwtSecret: []byte(secret),
	}
}

func (s *Service) VerifyAccessToken(tokenString string) (*UserContext, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)

	if userID == "" {
		return nil, fmt.Errorf("missing subject")
	}

	return &UserContext{
		SupabaseUserID: userID,
		Email:          email,
	}, nil
}