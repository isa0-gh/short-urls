package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/isa0-gh/short-urls/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *db.Queries
	secret  []byte
	exp     time.Duration
}

func NewAuthService(queries *db.Queries, secret []byte, exp time.Duration) *AuthService {
	return &AuthService{
		queries: queries,
		secret:  secret,
		exp:     exp,
	}
}

type CustomClaims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

func (a *AuthService) CheckPassword(
	ctx context.Context, password string, email string,
) (string, error) {
	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	return user.ID.String(), nil
}

func (a *AuthService) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(a.exp)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *AuthService) VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Prevent algorithm confusion / key substitution attacks (RFC 8725).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
