package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/orderfile/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderFileRepository interface {
	Create(ctx context.Context, file *model.OrderFile) error
	FindByOrderID(ctx context.Context, orderID string) ([]model.OrderFile, error)
	Delete(ctx context.Context, id string) error
}

type orderFileRepository struct {
	db *pgxpool.Pool
}

func NewOrderFileRepository(db *pgxpool.Pool) OrderFileRepository {
	return &orderFileRepository{db: db}
}

func (r *orderFileRepository) Create(ctx context.Context, file *model.OrderFile) error {
	query := `
		INSERT INTO order_files (
			order_id,
			uploaded_by,
			file_name,
			file_url,
			file_type,
			file_size
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		file.OrderID,
		file.UploadedBy,
		file.FileName,
		file.FileURL,
		file.FileType,
		file.FileSize,
	).Scan(
		&file.ID,
		&file.CreatedAt,
	)
}

func (r *orderFileRepository) FindByOrderID(ctx context.Context, orderID string) ([]model.OrderFile, error) {
	query := `
		SELECT id, order_id, uploaded_by, file_name, file_url, file_type, file_size, created_at
		FROM order_files
		WHERE order_id = $1::uuid
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []model.OrderFile{}

	for rows.Next() {
		var file model.OrderFile

		err := rows.Scan(
			&file.ID,
			&file.OrderID,
			&file.UploadedBy,
			&file.FileName,
			&file.FileURL,
			&file.FileType,
			&file.FileSize,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, rows.Err()
}

func (r *orderFileRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM order_files
		WHERE id = $1::uuid
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}