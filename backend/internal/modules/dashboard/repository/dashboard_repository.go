package repository

import (
	"context"

	"nexusweb-market/backend/internal/modules/dashboard/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	query := `
	SELECT
		(SELECT COUNT(*) FROM users) AS total_customers,
		(SELECT COUNT(*) FROM service_orders) AS total_orders,
		(SELECT COUNT(*) FROM services) AS total_services,
		(SELECT COUNT(*) FROM service_packages) AS total_packages,

		(SELECT COUNT(*) FROM service_orders WHERE status = 'PENDING') AS pending_orders,
		(SELECT COUNT(*) FROM service_orders WHERE status = 'IN_PROGRESS') AS in_progress_orders,
		(SELECT COUNT(*) FROM service_orders WHERE status = 'COMPLETED') AS completed_orders,

		(SELECT COUNT(*) FROM payments WHERE payment_status = 'PENDING') AS pending_payments,
		(SELECT COUNT(*) FROM payments WHERE payment_status = 'PAID') AS paid_payments,

		COALESCE((SELECT SUM(amount) FROM payments WHERE payment_status = 'PAID'), 0) AS total_revenue
	`

	var stats dto.DashboardStatsResponse

	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalCustomers,
		&stats.TotalOrders,
		&stats.TotalServices,
		&stats.TotalPackages,
		&stats.PendingOrders,
		&stats.InProgressOrders,
		&stats.CompletedOrders,
		&stats.PendingPayments,
		&stats.PaidPayments,
		&stats.TotalRevenue,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}