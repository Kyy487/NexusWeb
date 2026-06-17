package repository

import (
	"context"
	"encoding/json"

	"nexusweb-market/backend/internal/modules/package/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PackageRepository interface {
	FindAll(ctx context.Context) ([]model.Package, error)
	FindByID(ctx context.Context, id string) (*model.Package, error)
	Create(ctx context.Context, pkg *model.Package) error
	Update(ctx context.Context, pkg *model.Package) error
	Delete(ctx context.Context, id string) error
}

type packageRepository struct {
	db *pgxpool.Pool
}

func NewPackageRepository(db *pgxpool.Pool) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) FindAll(ctx context.Context) ([]model.Package, error) {
	query := `
		SELECT 
			sp.id,
			sp.service_id,
			s.name,
			sp.name,
			COALESCE(sp.description, ''),
			sp.price,
			sp.revision_count,
			sp.delivery_days,
			COALESCE(sp.features, '[]'::jsonb),
			sp.status,
			sp.created_at,
			sp.updated_at
		FROM service_packages sp
		JOIN services s ON s.id = sp.service_id
		ORDER BY sp.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []model.Package

	for rows.Next() {
		var pkg model.Package
		var featuresBytes []byte

		err := rows.Scan(
			&pkg.ID,
			&pkg.ServiceID,
			&pkg.ServiceName,
			&pkg.Name,
			&pkg.Description,
			&pkg.Price,
			&pkg.RevisionCount,
			&pkg.DeliveryDays,
			&featuresBytes,
			&pkg.Status,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(featuresBytes, &pkg.Features)
		packages = append(packages, pkg)
	}

	return packages, rows.Err()
}

func (r *packageRepository) FindByID(ctx context.Context, id string) (*model.Package, error) {
	query := `
		SELECT 
			sp.id,
			sp.service_id,
			s.name,
			sp.name,
			COALESCE(sp.description, ''),
			sp.price,
			sp.revision_count,
			sp.delivery_days,
			COALESCE(sp.features, '[]'::jsonb),
			sp.status,
			sp.created_at,
			sp.updated_at
		FROM service_packages sp
		JOIN services s ON s.id = sp.service_id
		WHERE sp.id = $1
		LIMIT 1
	`

	var pkg model.Package
	var featuresBytes []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&pkg.ID,
		&pkg.ServiceID,
		&pkg.ServiceName,
		&pkg.Name,
		&pkg.Description,
		&pkg.Price,
		&pkg.RevisionCount,
		&pkg.DeliveryDays,
		&featuresBytes,
		&pkg.Status,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(featuresBytes, &pkg.Features)

	return &pkg, nil
}

func (r *packageRepository) Create(ctx context.Context, pkg *model.Package) error {
	featuresJSON, _ := json.Marshal(pkg.Features)

	query := `
		INSERT INTO service_packages (
			service_id,
			name,
			description,
			price,
			revision_count,
			delivery_days,
			features,
			status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		pkg.ServiceID,
		pkg.Name,
		pkg.Description,
		pkg.Price,
		pkg.RevisionCount,
		pkg.DeliveryDays,
		featuresJSON,
		pkg.Status,
	).Scan(&pkg.ID, &pkg.CreatedAt, &pkg.UpdatedAt)
}

func (r *packageRepository) Update(ctx context.Context, pkg *model.Package) error {
	featuresJSON, _ := json.Marshal(pkg.Features)

	query := `
		UPDATE service_packages
		SET service_id = $1,
			name = $2,
			description = $3,
			price = $4,
			revision_count = $5,
			delivery_days = $6,
			features = $7,
			status = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`

	_, err := r.db.Exec(
		ctx,
		query,
		pkg.ServiceID,
		pkg.Name,
		pkg.Description,
		pkg.Price,
		pkg.RevisionCount,
		pkg.DeliveryDays,
		featuresJSON,
		pkg.Status,
		pkg.ID,
	)

	return err
}

func (r *packageRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM service_packages
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}