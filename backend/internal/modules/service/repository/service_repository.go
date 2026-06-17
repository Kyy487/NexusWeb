package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/service/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository interface {
	FindAll(ctx context.Context) ([]model.Service, error)
	FindByID(ctx context.Context, id string) (*model.Service, error)
	Create(ctx context.Context, service *model.Service) error
	Update(ctx context.Context, service *model.Service) error
	Delete(ctx context.Context, id string) error
}

type serviceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) FindAll(ctx context.Context) ([]model.Service, error) {
	query := `
		SELECT 
			s.id,
			s.category_id,
			c.name,
			s.name,
			s.slug,
			COALESCE(s.description, ''),
			s.base_price,
			s.estimated_days,
			s.status,
			s.created_at,
			s.updated_at
		FROM services s
		JOIN service_categories c ON c.id = s.category_id
		WHERE s.deleted_at IS NULL
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []model.Service{}

	for rows.Next() {
		var service model.Service

		err := rows.Scan(
			&service.ID,
			&service.CategoryID,
			&service.CategoryName,
			&service.Name,
			&service.Slug,
			&service.Description,
			&service.BasePrice,
			&service.EstimatedDays,
			&service.Status,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, rows.Err()
}

func (r *serviceRepository) FindByID(ctx context.Context, id string) (*model.Service, error) {
	service := &model.Service{}

	query := `
		SELECT 
			s.id,
			s.category_id,
			c.name,
			s.name,
			s.slug,
			COALESCE(s.description, ''),
			s.base_price,
			s.estimated_days,
			s.status,
			s.created_at,
			s.updated_at
		FROM services s
		JOIN service_categories c ON c.id = s.category_id
		WHERE s.id = $1
		AND s.deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.CategoryID,
		&service.CategoryName,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.BasePrice,
		&service.EstimatedDays,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *serviceRepository) Create(ctx context.Context, service *model.Service) error {
	query := `
		INSERT INTO services (
			category_id,
			name,
			slug,
			description,
			base_price,
			estimated_days,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		service.CategoryID,
		service.Name,
		service.Slug,
		service.Description,
		service.BasePrice,
		service.EstimatedDays,
		service.Status,
	).Scan(
		&service.ID,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
}

func (r *serviceRepository) Update(ctx context.Context, service *model.Service) error {
	query := `
		UPDATE services
		SET category_id = $1,
			name = $2,
			slug = $3,
			description = $4,
			base_price = $5,
			estimated_days = $6,
			status = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(
		ctx,
		query,
		service.CategoryID,
		service.Name,
		service.Slug,
		service.Description,
		service.BasePrice,
		service.EstimatedDays,
		service.Status,
		service.ID,
	)

	return err
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE services
		SET deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}