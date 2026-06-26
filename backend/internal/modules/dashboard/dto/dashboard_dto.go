package dto

type DashboardStatsResponse struct {
	TotalCustomers   int64   `json:"total_customers"`
	TotalOrders      int64   `json:"total_orders"`
	TotalServices    int64   `json:"total_services"`
	TotalPackages    int64   `json:"total_packages"`
	PendingOrders    int64   `json:"pending_orders"`
	InProgressOrders int64   `json:"in_progress_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	PendingPayments  int64   `json:"pending_payments"`
	PaidPayments     int64   `json:"paid_payments"`
	TotalRevenue     float64 `json:"total_revenue"`
}