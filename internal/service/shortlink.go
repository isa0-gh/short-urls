package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

const shortLinkCacheTTL = 5 * time.Minute

type cachedShortLinks struct {
	Limit int32     `json:"limit"`
	Links []db.Link `json:"links"`
}

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

	link, err := s.queries.CreateShortLink(ctx,
		db.CreateShortLinkParams{
			OwnerID:        ownerID,
			ShortCode:      shortCode,
			DestinationUrl: destinationUrl,
		})
	if err != nil {
		return db.Link{}, err
	}
	if s.redis != nil {
		_ = s.redis.Set(ctx, destinationCacheKey(link.ShortCode), link.DestinationUrl, shortLinkCacheTTL).Err()
		s.invalidateListCache(ctx, ownerID.String())
	}
	return link, nil
}

func (s *ShortLinkService) GetDestinationUrl(
	ctx context.Context,
	shortCode string,
) (string, error) {
	cacheKey := destinationCacheKey(shortCode)
	if s.redis != nil {
		if destination, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			return destination, nil
		}
	}

	destination, err := s.queries.GetDestinationURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if s.redis != nil {
		_ = s.redis.Set(ctx, cacheKey, destination, shortLinkCacheTTL).Err()
	}
	return destination, nil
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

	shortCode, err := s.queries.UpdateDestination(ctx, db.UpdateDestinationParams{
		DestinationUrl: destinationUrl,
		OwnerID:        ownerID,
		ID:             linkID,
	})
	if err != nil {
		return err
	}
	if s.redis != nil {
		_ = s.redis.Del(ctx, destinationCacheKey(shortCode)).Err()
		s.invalidateListCache(ctx, ownerID.String())
	}
	return nil
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

	shortCode, err := s.queries.DeleteShortLink(ctx, db.DeleteShortLinkParams{
		OwnerID: ownerID,
		ID:      linkID,
	})
	if err != nil {
		return err
	}
	if s.redis != nil {
		_ = s.redis.Del(ctx, destinationCacheKey(shortCode)).Err()
		s.invalidateListCache(ctx, ownerID.String())
	}
	return nil
}

func (s *ShortLinkService) ListShortLinks(
	ctx context.Context,
	userID string, limit int32,
) ([]db.Link, error) {
	var ownerID pgtype.UUID
	if err := ownerID.Scan(userID); err != nil {
		return []db.Link{}, err
	}
	cacheKey := shortLinkListCacheKey(ownerID.String())
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
			var entry cachedShortLinks
			if json.Unmarshal(cached, &entry) == nil && entry.Limit == limit {
				return entry.Links, nil
			}
		}
	}

	links, err := s.queries.ListShortLinks(ctx,
		db.ListShortLinksParams{OwnerID: ownerID, Limit: limit})
	if err == sql.ErrNoRows {
		return links, nil
	}
	if err != nil {
		return nil, err
	}
	if s.redis != nil {
		if cached, err := json.Marshal(cachedShortLinks{Limit: limit, Links: links}); err == nil {
			_ = s.redis.Set(ctx, cacheKey, cached, shortLinkCacheTTL).Err()
		}
	}
	return links, nil
}

func destinationCacheKey(shortCode string) string {
	return "shortlink:destination:" + shortCode
}

func shortLinkListCacheKey(userID string) string {
	return "shortlink:list:" + userID
}

func (s *ShortLinkService) invalidateListCache(ctx context.Context, userID string) {
	_ = s.redis.Del(ctx, shortLinkListCacheKey(userID)).Err()
}
