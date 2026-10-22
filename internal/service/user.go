package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/isa0-gh/short-urls/internal/db"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{
		queries: queries,
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
	return s.queries.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(
	ctx context.Context, id string,
) (db.User, error) {
	var userID pgtype.UUID
	if err := userID.Scan(id); err != nil {
		return db.User{}, err
	}
	return s.queries.GetUserByID(ctx, userID)
}
