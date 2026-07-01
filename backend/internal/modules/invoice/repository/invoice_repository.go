package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/invoice/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceRepository interface {
	FindAll(ctx context.Context) ([]model.Invoice, error)
	FindByID(ctx context.Context, id string) (*model.Invoice, error)
	FindByOrderID(ctx context.Context, orderID string) (*model.Invoice, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]model.Invoice, error)
	Create(ctx context.Context, invoice *model.Invoice) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetOrderAmount(ctx context.Context, orderID string) (float64, error)
	GetCustomerID(ctx context.Context, id string) (string, error)
	GetCustomerIDByOrderID(ctx context.Context, orderID string) (string, error)
}

type invoiceRepository struct {
	db *pgxpool.Pool
}

func NewInvoiceRepository(db *pgxpool.Pool) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) FindAll(ctx context.Context) ([]model.Invoice, error) {
	query := `
		SELECT
			i.id,
			i.order_id,
			so.order_number,
			i.invoice_number,
			i.subtotal,
			i.discount,
			i.tax,
			i.total_amount,
			i.status,
			i.due_date,
			i.created_at,
			i.updated_at
		FROM invoices i
		JOIN service_orders so ON so.id = i.order_id
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []model.Invoice{}

	for rows.Next() {
		var invoice model.Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.OrderID,
			&invoice.OrderNumber,
			&invoice.InvoiceNumber,
			&invoice.Subtotal,
			&invoice.Discount,
			&invoice.Tax,
			&invoice.TotalAmount,
			&invoice.Status,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	return invoices, rows.Err()
}

func (r *invoiceRepository) FindByID(ctx context.Context, id string) (*model.Invoice, error) {
	query := `
		SELECT
			i.id,
			i.order_id,
			so.order_number,
			i.invoice_number,
			i.subtotal,
			i.discount,
			i.tax,
			i.total_amount,
			i.status,
			i.due_date,
			i.created_at,
			i.updated_at
		FROM invoices i
		JOIN service_orders so ON so.id = i.order_id
		WHERE i.id = $1
		LIMIT 1
	`

	var invoice model.Invoice

	err := r.db.QueryRow(ctx, query, id).Scan(
		&invoice.ID,
		&invoice.OrderID,
		&invoice.OrderNumber,
		&invoice.InvoiceNumber,
		&invoice.Subtotal,
		&invoice.Discount,
		&invoice.Tax,
		&invoice.TotalAmount,
		&invoice.Status,
		&invoice.DueDate,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) FindByOrderID(ctx context.Context, orderID string) (*model.Invoice, error) {
	query := `
		SELECT
			i.id,
			i.order_id,
			so.order_number,
			i.invoice_number,
			i.subtotal,
			i.discount,
			i.tax,
			i.total_amount,
			i.status,
			i.due_date,
			i.created_at,
			i.updated_at
		FROM invoices i
		JOIN service_orders so ON so.id = i.order_id
		WHERE i.order_id = $1
		LIMIT 1
	`

	var invoice model.Invoice

	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&invoice.ID,
		&invoice.OrderID,
		&invoice.OrderNumber,
		&invoice.InvoiceNumber,
		&invoice.Subtotal,
		&invoice.Discount,
		&invoice.Tax,
		&invoice.TotalAmount,
		&invoice.Status,
		&invoice.DueDate,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) FindByCustomerID(ctx context.Context, customerID string) ([]model.Invoice, error) {
	query := `
		SELECT
			i.id,
			i.order_id,
			so.order_number,
			i.invoice_number,
			i.subtotal,
			i.discount,
			i.tax,
			i.total_amount,
			i.status,
			i.due_date,
			i.created_at,
			i.updated_at
		FROM invoices i
		JOIN service_orders so ON so.id = i.order_id
		WHERE so.customer_id = $1
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []model.Invoice{}

	for rows.Next() {
		var invoice model.Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.OrderID,
			&invoice.OrderNumber,
			&invoice.InvoiceNumber,
			&invoice.Subtotal,
			&invoice.Discount,
			&invoice.Tax,
			&invoice.TotalAmount,
			&invoice.Status,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	return invoices, rows.Err()
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *model.Invoice) error {
	query := `
		INSERT INTO invoices (
			order_id,
			invoice_number,
			subtotal,
			discount,
			tax,
			total_amount,
			status,
			due_date
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		invoice.OrderID,
		invoice.InvoiceNumber,
		invoice.Subtotal,
		invoice.Discount,
		invoice.Tax,
		invoice.TotalAmount,
		invoice.Status,
		invoice.DueDate,
	).Scan(
		&invoice.ID,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE invoices
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *invoiceRepository) GetOrderAmount(ctx context.Context, orderID string) (float64, error) {
	query := `
		SELECT total_price
		FROM service_orders
		WHERE id = $1
		LIMIT 1
	`

	var amount float64
	err := r.db.QueryRow(ctx, query, orderID).Scan(&amount)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (r *invoiceRepository) GetCustomerID(ctx context.Context, id string) (string, error) {
	query := `
		SELECT so.customer_id
		FROM invoices i
		JOIN service_orders so ON so.id = i.order_id
		WHERE i.id = $1
		LIMIT 1
	`
	var customerID string
	err := r.db.QueryRow(ctx, query, id).Scan(&customerID)
	if err != nil {
		return "", err
	}
	return customerID, nil
}

func (r *invoiceRepository) GetCustomerIDByOrderID(ctx context.Context, orderID string) (string, error) {
	query := `
		SELECT customer_id
		FROM service_orders
		WHERE id = $1
		LIMIT 1
	`
	var customerID string
	err := r.db.QueryRow(ctx, query, orderID).Scan(&customerID)
	if err != nil {
		return "", err
	}
	return customerID, nil
}
