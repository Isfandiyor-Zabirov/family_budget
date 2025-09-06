// /internal/auth/jwt_service.go
package auth

import (
	"family_budget/internal/internal_config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateTokens(userID int, familyID int, roleID int) (accessToken string, refreshToken string, err error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*RefreshClaims, error)
}

type Claims struct {
	UserID   int `json:"user_id"`
	FamilyID int `json:"family_id"`
	RoleID   int `json:"role_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	accessSecretKey  string
	refreshSecretKey string
}

func NewJWTService() JWTService {
	return &jwtService{
		accessSecretKey:  internal_config.InternalConfigs.Application.AccessKey,
		refreshSecretKey: internal_config.InternalConfigs.Application.RefreshKey,
	}
}

func (s *jwtService) GenerateTokens(userID int, familyID int, roleID int) (accessToken string, refreshToken string, err error) {
	accessTimeout := time.Duration(internal_config.InternalConfigs.Application.AccessTknTimeout) * time.Millisecond
	accessExpirationTime := time.Now().Add(accessTimeout)
	accessClaims := &Claims{
		UserID:   userID,
		FamilyID: familyID,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accToken.SignedString([]byte(s.accessSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshTimeout := time.Duration(internal_config.InternalConfigs.Application.RefreshTknTimeout) * time.Millisecond
	refreshExpirationTime := time.Now().Add(refreshTimeout)
	refreshClaims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refToken.SignedString([]byte(s.refreshSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.accessSecretKey), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid access token")
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.refreshSecretKey), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid refresh token")
}