package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const sessionCacheTTL = 5 * time.Minute

type AuthService struct {
	queries     *db.Queries
	redisClient *redis.Client
	exp         time.Duration
}

func NewAuthService(queries *db.Queries, redisClient *redis.Client, exp time.Duration) *AuthService {
	return &AuthService{
		queries:     queries,
		redisClient: redisClient,
		exp:         exp,
	}
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

func (a *AuthService) GenerateToken(ctx context.Context, userID string) (string, error) {
	randomBytes := make([]byte, 25)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(randomBytes)

	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return "", err
	}

	sessionHash := sha256.Sum256([]byte(token))
	_, err := a.queries.CreateSession(ctx, db.CreateSessionParams{
		SessionHash: hex.EncodeToString(sessionHash[:]),
		ID:          ownerID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(a.exp),
			Valid: true,
		},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) VerifyToken(ctx context.Context, token string) (string, error) {
	if len(token) != 34 {
		return "", errors.New("invalid session token")
	}

	randomBytes, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(randomBytes) != 25 || base64.RawURLEncoding.EncodeToString(randomBytes) != token {
		return "", errors.New("invalid session token")
	}

	sessionHash := sha256.Sum256([]byte(token))
	sessionHashHex := hex.EncodeToString(sessionHash[:])
	cacheKey := "session:" + sessionHashHex
	if a.redisClient != nil {
		cachedOwnerID, err := a.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var ownerID pgtype.UUID
			if ownerID.Scan(cachedOwnerID) == nil {
				return ownerID.String(), nil
			}
		}
	}

	ownerID, err := a.queries.GetUserIDBySessionHash(ctx, sessionHashHex)
	if err != nil {
		return "", err
	}
	if a.redisClient != nil {
		_ = a.redisClient.Set(ctx, cacheKey, ownerID.String(), sessionCacheTTL).Err()
	}
	return ownerID.String(), nil
}
