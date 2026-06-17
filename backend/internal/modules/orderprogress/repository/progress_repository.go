package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/orderprogress/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository interface {
	FindByOrderID(ctx context.Context, orderID string) ([]model.Progress, error)
	FindByID(ctx context.Context, id string) (*model.Progress, error)
	Create(ctx context.Context, progress *model.Progress) error
	Update(ctx context.Context, progress *model.Progress) error
	Delete(ctx context.Context, id string) error
}

type progressRepository struct {
	db *pgxpool.Pool
}

func NewProgressRepository(db *pgxpool.Pool) ProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) FindByOrderID(ctx context.Context, orderID string) ([]model.Progress, error) {
	query := `
		SELECT
			id,
			order_id,
			title,
			COALESCE(description, ''),
			progress_percentage,
			created_by,
			created_at
		FROM order_progress
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progressList := []model.Progress{}

	for rows.Next() {
		var progress model.Progress

		err := rows.Scan(
			&progress.ID,
			&progress.OrderID,
			&progress.Title,
			&progress.Description,
			&progress.ProgressPercentage,
			&progress.CreatedBy,
			&progress.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		progressList = append(progressList, progress)
	}

	return progressList, rows.Err()
}

func (r *progressRepository) FindByID(ctx context.Context, id string) (*model.Progress, error) {
	query := `
		SELECT
			id,
			order_id,
			title,
			COALESCE(description, ''),
			progress_percentage,
			created_by,
			created_at
		FROM order_progress
		WHERE id = $1
		LIMIT 1
	`

	var progress model.Progress

	err := r.db.QueryRow(ctx, query, id).Scan(
		&progress.ID,
		&progress.OrderID,
		&progress.Title,
		&progress.Description,
		&progress.ProgressPercentage,
		&progress.CreatedBy,
		&progress.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *progressRepository) Create(ctx context.Context, progress *model.Progress) error {
	query := `
		INSERT INTO order_progress (
			order_id,
			title,
			description,
			progress_percentage,
			created_by
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		progress.OrderID,
		progress.Title,
		progress.Description,
		progress.ProgressPercentage,
		progress.CreatedBy,
	).Scan(
		&progress.ID,
		&progress.CreatedAt,
	)
}

func (r *progressRepository) Update(ctx context.Context, progress *model.Progress) error {
	query := `
		UPDATE order_progress
		SET
			title = $1,
			description = $2,
			progress_percentage = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		progress.Title,
		progress.Description,
		progress.ProgressPercentage,
		progress.ID,
	)

	return err
}

func (r *progressRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM order_progress
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}