package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/activitylog/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *model.ActivityLog) error
	FindAll(ctx context.Context) ([]model.ActivityLog, error)
	FindByUserID(ctx context.Context, userID string) ([]model.ActivityLog, error)
}

type activityLogRepository struct {
	db *pgxpool.Pool
}

func NewActivityLogRepository(db *pgxpool.Pool) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *model.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (
			user_id,
			module,
			action,
			description,
			ip_address
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		log.UserID,
		log.Module,
		log.Action,
		log.Description,
		log.IPAddress,
	).Scan(
		&log.ID,
		&log.CreatedAt,
	)
}

func (r *activityLogRepository) FindAll(ctx context.Context) ([]model.ActivityLog, error) {
	query := `
		SELECT id, user_id, module, action, description, ip_address, created_at
		FROM activity_logs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ActivityLog{}

	for rows.Next() {
		var log model.ActivityLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Module,
			&log.Action,
			&log.Description,
			&log.IPAddress,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func (r *activityLogRepository) FindByUserID(ctx context.Context, userID string) ([]model.ActivityLog, error) {
	query := `
		SELECT id, user_id, module, action, description, ip_address, created_at
		FROM activity_logs
		WHERE user_id = $1::uuid
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ActivityLog{}

	for rows.Next() {
		var log model.ActivityLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Module,
			&log.Action,
			&log.Description,
			&log.IPAddress,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}