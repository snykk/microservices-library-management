package jwt

import (
	"errors"
	"time"

	driJWT "github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(userId string, Role string, email string) (t string, err error)
	GenerateRefreshToken(userID string, Role string, email string) (t string, err error)
	ParseToken(tokenString string) (claims JwtCustomClaim, err error)
}

type JwtCustomClaim struct {
	UserID string
	Role   string
	Email  string
	driJWT.StandardClaims
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
		userID,
		Role,
		email,
		driJWT.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiredAccessToken).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := driJWT.NewWithClaims(driJWT.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

func (j *jwtService) GenerateRefreshToken(userID string, Role string, email string) (t string, err error) {
	claims := &JwtCustomClaim{
		userID,
		Role,
		email,
		driJWT.StandardClaims{
			ExpiresAt: time.Now().Add(j.expiredRefreshToken).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := driJWT.NewWithClaims(driJWT.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

// ParseToken parses and validates the JWT token, extracting the claims.
func (j *jwtService) ParseToken(tokenString string) (claims JwtCustomClaim, err error) {
	if token, err := driJWT.ParseWithClaims(tokenString, &claims, func(token *driJWT.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	}); err != nil || !token.Valid {
		return JwtCustomClaim{}, errors.New("token is not valid")
	}

	return
}
