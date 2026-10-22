package service

import (
	"context"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type TelemetryService struct {
	queries *db.Queries
}

func NewTelemetryService(queries *db.Queries) *TelemetryService {
	return &TelemetryService{
		queries: queries,
	}
}

func (s *TelemetryService) CreateTelemetry(
	ctx context.Context, linkID string,
	referrerURL pgtype.Text, userAgent pgtype.Text, ipAddress pgtype.Text,
) error {
	var shortCode pgtype.UUID
	if err := shortCode.Scan(linkID); err != nil {
		return err
	}
	return s.queries.CreateTelemetry(ctx, db.CreateTelemetryParams{
		LinkID:      shortCode,
		ReferrerUrl: referrerURL,
		UserAgent:   userAgent,
		IpAddress:   ipAddress,
	})
}
