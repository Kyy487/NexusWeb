package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindRoleByName(ctx context.Context, name string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) FindRoleByName(ctx context.Context, name string) (string, error) {
	var id string

	query := `SELECT id FROM roles WHERE name = $1 LIMIT 1`

	err := r.db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	query := `
		SELECT 
			u.id,
			u.role_id,
			r.name,
			u.name,
			u.email,
			u.password_hash,
			COALESCE(u.phone, ''),
			u.status,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.RoleName,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
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

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := `
		INSERT INTO users (role_id, name, email, password_hash, phone, status)
		VALUES ($1, $2, $3, $4, $5, 'ACTIVE')
		RETURNING id, role_id, name, email, COALESCE(phone, ''), status, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.RoleID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Phone,
	).Scan(
		&user.ID,
		&user.RoleID,
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