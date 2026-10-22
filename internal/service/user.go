package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/redis/go-redis/v9"
)

const userCacheTTL = 2 * time.Minute

type UserService struct {
	queries *db.Queries
	redis   *redis.Client
}

func NewUserService(queries *db.Queries, redis *redis.Client) *UserService {
	return &UserService{
		queries: queries,
		redis:   redis,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context, username string,
	email string, passwordHash string,
) (db.User, error) {
	return s.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
}

func (s *UserService) GetUserByEmail(
	ctx context.Context, email string,
) (db.User, error) {
	emailHash := sha256.Sum256([]byte(email))
	return s.getCachedUser(ctx, "user:email:"+hex.EncodeToString(emailHash[:]), func() (db.User, error) {
		return s.queries.GetUserByEmail(ctx, email)
	})
}

func (s *UserService) GetUserByID(
	ctx context.Context, id string,
) (db.User, error) {
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		return db.User{}, err
	}
	return s.getCachedUser(ctx, "user:id:"+userID.String(), func() (db.User, error) {
		return s.queries.GetUserByID(ctx, userID)
	})
}

func (s *UserService) getCachedUser(
	ctx context.Context,
	cacheKey string,
	load func() (db.User, error),
) (db.User, error) {
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
			var user db.User
			if json.Unmarshal(cached, &user) == nil {
				return user, nil
			}
		}
	}

	user, err := load()
	if err != nil {
		return db.User{}, err
	}
	if s.redis != nil {
		if cached, err := json.Marshal(user); err == nil {
			_ = s.redis.Set(ctx, cacheKey, cached, userCacheTTL).Err()
		}
	}
	return user, nil
}
