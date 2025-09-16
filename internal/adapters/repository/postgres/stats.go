package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/itsbaivab/url-shortener/internal/core/domain"
	_ "github.com/lib/pq"
)

type PostgresStatsRepository struct {
	db *sql.DB
}

func NewPostgresStatsRepository(db *sql.DB) *PostgresStatsRepository {
	return &PostgresStatsRepository{db: db}
}

func (r *PostgresStatsRepository) All(ctx context.Context) ([]domain.Stats, error) {
	query := `SELECT id, link_id, platform, created_at FROM stats ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	var stats []domain.Stats
	for rows.Next() {
		var stat domain.Stats
		err := rows.Scan(&stat.Id, &stat.LinkID, &stat.Platform, &stat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stat: %w", err)
		}
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return stats, nil
}

func (r *PostgresStatsRepository) Get(ctx context.Context, id string) (domain.Stats, error) {
	var stat domain.Stats
	query := `SELECT id, link_id, platform, created_at FROM stats WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stat.Id,
		&stat.LinkID,
		&stat.Platform,
		&stat.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.Stats{}, fmt.Errorf("stats not found")
	}
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}

	return stat, nil
}

func (r *PostgresStatsRepository) Create(ctx context.Context, stats domain.Stats) error {
	query := `INSERT INTO stats (id, link_id, platform, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, stats.Id, stats.LinkID, stats.Platform, stats.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create stats: %w", err)
	}

	return nil
}

func (r *PostgresStatsRepository) Delete(ctx context.Context, linkID string) error {
	query := `DELETE FROM stats WHERE link_id = $1`

	_, err := r.db.ExecContext(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed to delete stats: %w", err)
	}

	return nil
}

func (r *PostgresStatsRepository) GetStatsByLinkID(ctx context.Context, linkID string) ([]domain.Stats, error) {
	query := `SELECT id, link_id, platform, created_at FROM stats WHERE link_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats by link ID: %w", err)
	}
	defer rows.Close()

	var stats []domain.Stats
	for rows.Next() {
		var stat domain.Stats
		err := rows.Scan(&stat.Id, &stat.LinkID, &stat.Platform, &stat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stat: %w", err)
		}
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return stats, nil
}
