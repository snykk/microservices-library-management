package jwt

import (
	"auth_service/internal/constants"
	"errors"
	"time"

	golangJWT "github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId string, Role string, email string) (t string, err error)
	GenerateRefreshToken(userID string, Role string, email string) (t string, err error)
	ParseToken(tokenString string, expectedType string) (claims JwtCustomClaim, err error)
}

type JwtCustomClaim struct {
	UserID    string
	Role      string
	Email     string
	TokenType string
	golangJWT.RegisteredClaims
}

type jwtService struct {
	secretKey           string
	issuer              string
	expiredAccessToken  time.Duration
	expiredRefreshToken time.Duration
}

func NewJWTService(secretKey, issuer string, expiredAccessToken, expiredRefreshToken time.Duration) JWTService {
	return &jwtService{
		issuer:              issuer,
		secretKey:           secretKey,
		expiredAccessToken:  expiredAccessToken,
		expiredRefreshToken: expiredRefreshToken,
	}
}

// GenerateToken creates a new JWT token for authentication with short expiry (e.g., 15 minutes)
func (j *jwtService) GenerateToken(userID string, Role string, email string) (t string, err error) {
	claims := &JwtCustomClaim{
		UserID:    userID,
		Role:      Role,
		Email:     email,
		TokenType: constants.TokenAccess,
		RegisteredClaims: golangJWT.RegisteredClaims{
			ExpiresAt: golangJWT.NewNumericDate(time.Now().Add(j.expiredAccessToken)),
			Issuer:    j.issuer,
			IssuedAt:  golangJWT.NewNumericDate(time.Now()),
		},
	}
	token := golangJWT.NewWithClaims(golangJWT.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

func (j *jwtService) GenerateRefreshToken(userID string, Role string, email string) (t string, err error) {
	claims := &JwtCustomClaim{
		UserID:    userID,
		Role:      Role,
		Email:     email,
		TokenType: constants.TokenRefresh,
		RegisteredClaims: golangJWT.RegisteredClaims{
			ExpiresAt: golangJWT.NewNumericDate(time.Now().Add(j.expiredRefreshToken)),
			Issuer:    j.issuer,
			IssuedAt:  golangJWT.NewNumericDate(time.Now()),
		},
	}
	token := golangJWT.NewWithClaims(golangJWT.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

// ParseToken parses and validates the JWT token, extracting the claims.
func (j *jwtService) ParseToken(tokenString string, expectedType string) (claims JwtCustomClaim, err error) {
	token, err := golangJWT.ParseWithClaims(tokenString, &claims, func(token *golangJWT.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return JwtCustomClaim{}, errors.New("token is not valid")
	}

	if claims.TokenType != expectedType {
		return JwtCustomClaim{}, errors.New("invalid token type")
	}

	return claims, nil
}
