package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/order/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	FindAll(ctx context.Context) ([]model.Order, error)
	FindByID(ctx context.Context, id string) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetPackagePrice(ctx context.Context, packageID string) (float64, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll(ctx context.Context) ([]model.Order, error) {
	query := `
		SELECT
			so.id,
			so.customer_id,
			COALESCE(u.name, ''),
			so.service_id,
			s.name,
			so.package_id,
			sp.name,
			so.order_number,
			so.title,
			COALESCE(so.description, ''),
			so.deadline,
			so.total_price,
			so.status,
			so.created_at,
			so.updated_at,
			so.completed_at,
			so.cancelled_at
		FROM service_orders so
		JOIN users u ON u.id = so.customer_id
		JOIN services s ON s.id = so.service_id
		JOIN service_packages sp ON sp.id = so.package_id
		ORDER BY so.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.CustomerName,
			&order.ServiceID,
			&order.ServiceName,
			&order.PackageID,
			&order.PackageName,
			&order.OrderNumber,
			&order.Title,
			&order.Description,
			&order.Deadline,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CompletedAt,
			&order.CancelledAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	query := `
		SELECT
			so.id,
			so.customer_id,
			COALESCE(u.name, ''),
			so.service_id,
			s.name,
			so.package_id,
			sp.name,
			so.order_number,
			so.title,
			COALESCE(so.description, ''),
			so.deadline,
			so.total_price,
			so.status,
			so.created_at,
			so.updated_at,
			so.completed_at,
			so.cancelled_at
		FROM service_orders so
		JOIN users u ON u.id = so.customer_id
		JOIN services s ON s.id = so.service_id
		JOIN service_packages sp ON sp.id = so.package_id
		WHERE so.id = $1
		LIMIT 1
	`

	var order model.Order

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.CustomerName,
		&order.ServiceID,
		&order.ServiceName,
		&order.PackageID,
		&order.PackageName,
		&order.OrderNumber,
		&order.Title,
		&order.Description,
		&order.Deadline,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CompletedAt,
		&order.CancelledAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO service_orders (
			customer_id,
			service_id,
			package_id,
			order_number,
			title,
			description,
			deadline,
			total_price,
			status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		order.CustomerID,
		order.ServiceID,
		order.PackageID,
		order.OrderNumber,
		order.Title,
		order.Description,
		order.Deadline,
		order.TotalPrice,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *orderRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {

	query := `
		UPDATE service_orders
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)

	return err
}
func (r *orderRepository) GetPackagePrice(ctx context.Context, packageID string) (float64, error) {
	query := `
		SELECT price
		FROM service_packages
		WHERE id = $1
		LIMIT 1
	`

	var price float64
	err := r.db.QueryRow(ctx, query, packageID).Scan(&price)
	if err != nil {
		return 0, err
	}

	return price, nil
}