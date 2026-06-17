package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT 
			u.id,
			u.role_id,
			r.name,
			u.name,
			u.email,
			COALESCE(u.phone, ''),
			u.status,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
		AND u.deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.RoleName,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT 
			u.id,
			u.role_id,
			r.name,
			u.name,
			u.email,
			COALESCE(u.phone, ''),
			u.status,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}

	for rows.Next() {
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.RoleName,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}