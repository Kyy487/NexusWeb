package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/category/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id string) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(description, ''), status, created_at, updated_at
		FROM service_categories
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []model.Category{}

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Status,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*model.Category, error) {
	category := &model.Category{}

	query := `
		SELECT id, name, slug, COALESCE(description, ''), status, created_at, updated_at
		FROM service_categories
		WHERE id = $1
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO service_categories (name, slug, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.Status,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE service_categories
		SET name = $1,
			slug = $2,
			description = $3,
			status = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err := r.db.Exec(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.Status,
		category.ID,
	)

	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM service_categories
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}