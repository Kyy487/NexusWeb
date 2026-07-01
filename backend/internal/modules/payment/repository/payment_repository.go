package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/payment/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	FindAll(ctx context.Context) ([]model.Payment, error)
	FindByID(ctx context.Context, id string) (*model.Payment, error)
	FindByInvoiceID(ctx context.Context, invoiceID string) ([]model.Payment, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]model.Payment, error)
	Create(ctx context.Context, payment *model.Payment) error
	UpdateStatus(ctx context.Context, id string, paymentStatus string, verifiedBy *string) error
	GetWhatsAppData(ctx context.Context, paymentID string) (*model.Payment, error)
	UpdateInvoiceAndOrderAfterPayment(ctx context.Context, paymentID string) error
	GetCustomerIDByPaymentID(ctx context.Context, paymentID string) (string, error)
}

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) FindAll(ctx context.Context) ([]model.Payment, error) {
	query := `
		SELECT id, invoice_id, amount, payment_method, payment_status,
		       payment_proof_url, paid_at, verified_by, created_at, updated_at
		FROM payments
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []model.Payment{}

	for rows.Next() {
		var payment model.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.InvoiceID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.PaymentProofURL,
			&payment.PaidAt,
			&payment.VerifiedBy,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*model.Payment, error) {
	query := `
		SELECT id, invoice_id, amount, payment_method, payment_status,
		       payment_proof_url, paid_at, verified_by, created_at, updated_at
		FROM payments
		WHERE id = $1
		LIMIT 1
	`

	var payment model.Payment

	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.InvoiceID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.PaymentProofURL,
		&payment.PaidAt,
		&payment.VerifiedBy,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) FindByInvoiceID(ctx context.Context, invoiceID string) ([]model.Payment, error) {
	query := `
		SELECT id, invoice_id, amount, payment_method, payment_status,
		       payment_proof_url, paid_at, verified_by, created_at, updated_at
		FROM payments
		WHERE invoice_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []model.Payment{}

	for rows.Next() {
		var payment model.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.InvoiceID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.PaymentProofURL,
			&payment.PaidAt,
			&payment.VerifiedBy,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *paymentRepository) FindByCustomerID(ctx context.Context, customerID string) ([]model.Payment, error) {
	query := `
		SELECT p.id, p.invoice_id, p.amount, p.payment_method, p.payment_status,
		       p.payment_proof_url, p.paid_at, p.verified_by, p.created_at, p.updated_at
		FROM payments p
		JOIN invoices i ON i.id = p.invoice_id
		JOIN service_orders so ON so.id = i.order_id
		WHERE so.customer_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []model.Payment{}

	for rows.Next() {
		var payment model.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.InvoiceID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.PaymentProofURL,
			&payment.PaidAt,
			&payment.VerifiedBy,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (
			invoice_id,
			amount,
			payment_method,
			payment_status,
			payment_proof_url
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		payment.InvoiceID,
		payment.Amount,
		payment.PaymentMethod,
		payment.PaymentStatus,
		payment.PaymentProofURL,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, paymentStatus string, verifiedBy *string) error {
	query := `
		UPDATE payments
		SET payment_status = $1::varchar,
		    verified_by = COALESCE($2::uuid, verified_by),
		    paid_at = CASE 
		        WHEN $1::varchar = 'PAID' THEN CURRENT_TIMESTAMP 
		        ELSE paid_at 
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3::uuid
	`

	_, err := r.db.Exec(ctx, query, paymentStatus, verifiedBy, id)
	return err
}
func (r *paymentRepository) GetWhatsAppData(ctx context.Context, paymentID string) (*model.Payment, error) {
	query := `
		SELECT id, invoice_id, amount, payment_method, payment_status,
		       payment_proof_url, paid_at, verified_by, created_at, updated_at
		FROM payments
		WHERE id = $1::uuid
		LIMIT 1
	`

	var payment model.Payment

	err := r.db.QueryRow(ctx, query, paymentID).Scan(
		&payment.ID,
		&payment.InvoiceID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.PaymentProofURL,
		&payment.PaidAt,
		&payment.VerifiedBy,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}
func (r *paymentRepository) UpdateInvoiceAndOrderAfterPayment(ctx context.Context, paymentID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var invoiceID string
	var orderID string
	var verifiedBy *string

	err = tx.QueryRow(ctx, `
		SELECT p.invoice_id, i.order_id, p.verified_by
		FROM payments p
		JOIN invoices i ON i.id = p.invoice_id
		WHERE p.id = $1::uuid
		LIMIT 1
	`, paymentID).Scan(&invoiceID, &orderID, &verifiedBy)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE invoices
		SET status = 'PAID',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1::uuid
	`, invoiceID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE service_orders
		SET status = 'IN_PROGRESS',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1::uuid
	`, orderID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO order_progress (
			order_id,
			title,
			description,
			progress_percentage,
			created_by
		)
		SELECT
			$1::uuid,
			'Pembayaran Diverifikasi',
			'Pembayaran customer sudah diverifikasi dan project mulai diproses.',
			10,
			$2::uuid
		WHERE NOT EXISTS (
			SELECT 1
			FROM order_progress
			WHERE order_id = $1::uuid
			AND progress_percentage = 10
			AND title = 'Pembayaran Diverifikasi'
		)
	`, orderID, verifiedBy)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO activity_logs (
			user_id,
			module,
			action,
			description,
			ip_address
		)
		SELECT
			$1::uuid,
			'PAYMENT',
			'VERIFY',
			'Payment verified, invoice paid, order moved to IN_PROGRESS, and progress 10% created.',
			NULL
		WHERE NOT EXISTS (
			SELECT 1
			FROM activity_logs
			WHERE user_id = $1::uuid
			AND module = 'PAYMENT'
			AND action = 'VERIFY'
			AND description = 'Payment verified, invoice paid, order moved to IN_PROGRESS, and progress 10% created.'
		)
	`, verifiedBy)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *paymentRepository) GetCustomerIDByPaymentID(ctx context.Context, paymentID string) (string, error) {
	query := `
		SELECT so.customer_id
		FROM payments p
		JOIN invoices i ON i.id = p.invoice_id
		JOIN service_orders so ON so.id = i.order_id
		WHERE p.id = $1::uuid
		LIMIT 1
	`
	var customerID string
	err := r.db.QueryRow(ctx, query, paymentID).Scan(&customerID)
	if err != nil {
		return "", err
	}
	return customerID, nil
}
