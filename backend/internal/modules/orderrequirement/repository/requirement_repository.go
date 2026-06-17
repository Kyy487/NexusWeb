package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/orderrequirement/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequirementRepository interface {
	FindByOrderID(ctx context.Context, orderID string) ([]model.Requirement, error)
	FindByID(ctx context.Context, id string) (*model.Requirement, error)
	Create(ctx context.Context, requirement *model.Requirement) error
	Update(ctx context.Context, requirement *model.Requirement) error
	Delete(ctx context.Context, id string) error
}

type requirementRepository struct {
	db *pgxpool.Pool
}

func NewRequirementRepository(db *pgxpool.Pool) RequirementRepository {
	return &requirementRepository{db: db}
}

func (r *requirementRepository) FindByOrderID(ctx context.Context, orderID string) ([]model.Requirement, error) {
	query := `
		SELECT
			id,
			order_id,
			question,
			COALESCE(answer, ''),
			created_at
		FROM order_requirements
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requirements := []model.Requirement{}

	for rows.Next() {
		var requirement model.Requirement

		err := rows.Scan(
			&requirement.ID,
			&requirement.OrderID,
			&requirement.Question,
			&requirement.Answer,
			&requirement.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		requirements = append(requirements, requirement)
	}

	return requirements, rows.Err()
}

func (r *requirementRepository) FindByID(ctx context.Context, id string) (*model.Requirement, error) {
	query := `
		SELECT
			id,
			order_id,
			question,
			COALESCE(answer, ''),
			created_at
		FROM order_requirements
		WHERE id = $1
		LIMIT 1
	`

	var requirement model.Requirement

	err := r.db.QueryRow(ctx, query, id).Scan(
		&requirement.ID,
		&requirement.OrderID,
		&requirement.Question,
		&requirement.Answer,
		&requirement.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &requirement, nil
}

func (r *requirementRepository) Create(ctx context.Context, requirement *model.Requirement) error {
	query := `
		INSERT INTO order_requirements (
			order_id,
			question,
			answer
		)
		VALUES ($1,$2,$3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		requirement.OrderID,
		requirement.Question,
		requirement.Answer,
	).Scan(
		&requirement.ID,
		&requirement.CreatedAt,
	)
}

func (r *requirementRepository) Update(ctx context.Context, requirement *model.Requirement) error {
	query := `
		UPDATE order_requirements
		SET
			question = $1,
			answer = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(
		ctx,
		query,
		requirement.Question,
		requirement.Answer,
		requirement.ID,
	)

	return err
}

func (r *requirementRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM order_requirements
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}