package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/itsbaivab/url-shortener/internal/core/domain"
	_ "github.com/lib/pq"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostgresLinkRepository(db *sql.DB) *PostgresLinkRepository {
	return &PostgresLinkRepository{db: db}
}

func (r *PostgresLinkRepository) All(ctx context.Context) ([]domain.Link, error) {
	query := `SELECT id, original_url, created_at FROM links ORDER BY created_at DESC LIMIT 100`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	var links []domain.Link
	for rows.Next() {
		var link domain.Link
		err := rows.Scan(&link.Id, &link.OriginalURL, &link.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return links, nil
}

func (r *PostgresLinkRepository) Get(ctx context.Context, id string) (domain.Link, error) {
	var link domain.Link
	query := `SELECT id, original_url, created_at FROM links WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.Id,
		&link.OriginalURL,
		&link.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.Link{}, fmt.Errorf("link not found")
	}
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get link: %w", err)
	}

	return link, nil
}

func (r *PostgresLinkRepository) Create(ctx context.Context, link domain.Link) error {
	query := `INSERT INTO links (id, original_url, created_at) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, link.Id, link.OriginalURL, link.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	return nil
}

func (r *PostgresLinkRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM links WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("link not found")
	}

	return nil
}
