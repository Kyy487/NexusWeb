package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nexusweb-market/backend/internal/modules/order/dto"
	"nexusweb-market/backend/internal/modules/order/model"
	"nexusweb-market/backend/internal/modules/order/repository"
)

type OrderService interface {
	GetAll(ctx context.Context) ([]dto.OrderResponse, error)
	GetByID(ctx context.Context, id string) (*dto.OrderResponse, error)
	Create(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	UpdateStatus(ctx context.Context, id string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetAll(ctx context.Context) ([]dto.OrderResponse, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, toOrderResponse(order))
	}

	return responses, nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*dto.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order id is required")
	}

	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toOrderResponse(*order)
	return &response, nil
}

func (s *orderService) Create(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	price, err := s.repo.GetPackagePrice(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	var deadline *time.Time
	if req.Deadline != "" {
		parsedDeadline, err := time.Parse("2006-01-02", req.Deadline)
		if err != nil {
			return nil, errors.New("deadline format must be YYYY-MM-DD")
		}
		deadline = &parsedDeadline
	}

	order := &model.Order{
		CustomerID:  req.CustomerID,
		ServiceID:   req.ServiceID,
		PackageID:   req.PackageID,
		OrderNumber: generateOrderNumber(),
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
		TotalPrice:  price,
		Status:      "PENDING",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	fullOrder, err := s.repo.FindByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	response := toOrderResponse(*fullOrder)
	return &response, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, id string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order id is required")
	}

	if !isValidOrderStatus(req.Status) {
		return nil, errors.New("invalid order status")
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}

	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toOrderResponse(*order)
	return &response, nil
}

func toOrderResponse(order model.Order) dto.OrderResponse {
	var deadline *string
	if order.Deadline != nil {
		formatted := order.Deadline.Format("2006-01-02")
		deadline = &formatted
	}

	return dto.OrderResponse{
		ID:           order.ID,
		CustomerID:   order.CustomerID,
		CustomerName: order.CustomerName,
		ServiceID:    order.ServiceID,
		ServiceName:  order.ServiceName,
		PackageID:    order.PackageID,
		PackageName:  order.PackageName,
		OrderNumber:  order.OrderNumber,
		Title:        order.Title,
		Description:  order.Description,
		Deadline:     deadline,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

func isValidOrderStatus(status string) bool {
	validStatuses := map[string]bool{
		"PENDING":     true,
		"PAID":        true,
		"IN_PROGRESS": true,
		"REVISION":    true,
		"COMPLETED":   true,
		"CANCELLED":   true,
	}

	return validStatuses[status]
}