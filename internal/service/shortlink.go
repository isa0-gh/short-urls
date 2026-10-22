package service

import (
	"context"
	"database/sql"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

type ShortLinkService struct {
	queries *db.Queries
	redis   *redis.Client
}

func NewShortLinkService(queries *db.Queries, redis *redis.Client) *ShortLinkService {
	return &ShortLinkService{
		queries: queries,
		redis:   redis,
	}
}

func (s *ShortLinkService) CreateShortLink(
	ctx context.Context,
	userID string,
	shortCode string,
	destinationUrl string,
) (db.Link, error) {
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return db.Link{}, err
	}

	return s.queries.CreateShortLink(ctx,
		db.CreateShortLinkParams{
			OwnerID:        ownerID,
			ShortCode:      shortCode,
			DestinationUrl: destinationUrl,
		})
}

func (s *ShortLinkService) GetDestinationUrl(
	ctx context.Context,
	shortCode string,
) (string, error) {
	return s.queries.GetDestinationURL(ctx, shortCode)
}

func (s *ShortLinkService) UpdateDestination(
	ctx context.Context,
	userID string,
	id string,
	destinationUrl string,
) error {
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return err
	}

	var linkID pgtype.UUID
	if err := linkID.Scan(id); err != nil {
		return err
	}

	return s.queries.UpdateDestination(ctx, db.UpdateDestinationParams{
		DestinationUrl: destinationUrl,
		OwnerID:        ownerID,
		ID:             linkID,
	})
}

func (s *ShortLinkService) DeleteShortLink(
	ctx context.Context,
	userID string,
	id string,
) error {
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return err
	}

	var linkID pgtype.UUID
	if err := linkID.Scan(id); err != nil {
		return err
	}

	return s.queries.DeleteShortLink(ctx, db.DeleteShortLinkParams{
		OwnerID: ownerID,
		ID:      linkID,
	})
}

func (s *ShortLinkService) ListShortLinks(
	ctx context.Context,
	userID string, limit int32,
) ([]db.Link, error) {
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return []db.Link{}, err
	}
	links, err := s.queries.ListShortLinks(ctx,
		db.ListShortLinksParams{OwnerID: ownerID, Limit: limit})
	if err == sql.ErrNoRows {
		return links, nil
	}
	return links, err
}
