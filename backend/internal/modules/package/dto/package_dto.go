package dto

type PackageResponse struct {
	ID            string   `json:"id"`
	ServiceID     string   `json:"service_id"`
	ServiceName   string   `json:"service_name"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	RevisionCount int      `json:"revision_count"`
	DeliveryDays  int      `json:"delivery_days"`
	Features      []string `json:"features"`
	Status        string   `json:"status"`
}

type CreatePackageRequest struct {
	ServiceID     string   `json:"service_id" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required"`
	RevisionCount int      `json:"revision_count"`
	DeliveryDays  int      `json:"delivery_days"`
	Features      []string `json:"features"`
	Status        string   `json:"status"`
}

type UpdatePackageRequest struct {
	ServiceID     string   `json:"service_id" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required"`
	RevisionCount int      `json:"revision_count"`
	DeliveryDays  int      `json:"delivery_days"`
	Features      []string `json:"features"`
	Status        string   `json:"status"`
}