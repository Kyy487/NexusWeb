package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nexusweb-market/backend/internal/modules/invoice/dto"
	"nexusweb-market/backend/internal/modules/invoice/model"
	"nexusweb-market/backend/internal/modules/invoice/repository"
)

type InvoiceService interface {
	GetAll(ctx context.Context, userID string, role string) ([]dto.InvoiceResponse, error)
	GetByID(ctx context.Context, id string, userID string, role string) (*dto.InvoiceResponse, error)
	GetByOrderID(ctx context.Context, orderID string, userID string, role string) (*dto.InvoiceResponse, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]dto.InvoiceResponse, error)
	Create(ctx context.Context, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error)
	UpdateStatus(ctx context.Context, id string, req dto.UpdateInvoiceStatusRequest) (*dto.InvoiceResponse, error)
}

type invoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) GetAll(ctx context.Context, userID string, role string) ([]dto.InvoiceResponse, error) {
	var invoices []model.Invoice
	var err error

	if role == "CUSTOMER" {
		invoices, err = s.repo.FindByCustomerID(ctx, userID)
	} else {
		invoices, err = s.repo.FindAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	responses := make([]dto.InvoiceResponse, 0, len(invoices))
	for _, invoice := range invoices {
		responses = append(responses, toInvoiceResponse(invoice))
	}

	return responses, nil
}

func (s *invoiceService) GetByID(ctx context.Context, id string, userID string, role string) (*dto.InvoiceResponse, error) {
	if id == "" {
		return nil, errors.New("invoice id is required")
	}

	invoice, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == "CUSTOMER" {
		ownerID, err := s.repo.GetCustomerID(ctx, id)
		if err != nil || ownerID != userID {
			return nil, errors.New("you do not have permission to access this invoice")
		}
	}

	response := toInvoiceResponse(*invoice)
	return &response, nil
}

func (s *invoiceService) GetByOrderID(ctx context.Context, orderID string, userID string, role string) (*dto.InvoiceResponse, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	if role == "CUSTOMER" {
		ownerID, err := s.repo.GetCustomerIDByOrderID(ctx, orderID)
		if err != nil || ownerID != userID {
			return nil, errors.New("you do not have permission to access this invoice")
		}
	}

	invoice, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	response := toInvoiceResponse(*invoice)
	return &response, nil
}

func (s *invoiceService) GetByCustomerID(ctx context.Context, customerID string) ([]dto.InvoiceResponse, error) {
	if customerID == "" {
		return nil, errors.New("customer id is required")
	}

	invoices, err := s.repo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.InvoiceResponse, 0, len(invoices))
	for _, invoice := range invoices {
		responses = append(responses, toInvoiceResponse(invoice))
	}

	return responses, nil
}

func (s *invoiceService) Create(ctx context.Context, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	subtotal, err := s.repo.GetOrderAmount(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsedDueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return nil, errors.New("due_date format must be YYYY-MM-DD")
		}
		dueDate = &parsedDueDate
	}

	totalAmount := subtotal - req.Discount + req.Tax
	if totalAmount < 0 {
		return nil, errors.New("total amount cannot be negative")
	}

	invoice := &model.Invoice{
		OrderID:       req.OrderID,
		InvoiceNumber: generateInvoiceNumber(),
		Subtotal:      subtotal,
		Discount:      req.Discount,
		Tax:           req.Tax,
		TotalAmount:   totalAmount,
		Status:        "UNPAID",
		DueDate:       dueDate,
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	fullInvoice, err := s.repo.FindByID(ctx, invoice.ID)
	if err != nil {
		return nil, err
	}

	response := toInvoiceResponse(*fullInvoice)
	return &response, nil
}

func (s *invoiceService) UpdateStatus(ctx context.Context, id string, req dto.UpdateInvoiceStatusRequest) (*dto.InvoiceResponse, error) {
	if id == "" {
		return nil, errors.New("invoice id is required")
	}

	if !isValidInvoiceStatus(req.Status) {
		return nil, errors.New("invalid invoice status")
	}

	if err := s.repo.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}

	invoice, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toInvoiceResponse(*invoice)
	return &response, nil
}

func toInvoiceResponse(invoice model.Invoice) dto.InvoiceResponse {
	var dueDate *string
	if invoice.DueDate != nil {
		formatted := invoice.DueDate.Format("2006-01-02")
		dueDate = &formatted
	}

	return dto.InvoiceResponse{
		ID:            invoice.ID,
		OrderID:       invoice.OrderID,
		OrderNumber:   invoice.OrderNumber,
		InvoiceNumber: invoice.InvoiceNumber,
		Subtotal:      invoice.Subtotal,
		Discount:      invoice.Discount,
		Tax:           invoice.Tax,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		DueDate:       dueDate,
	}
}

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", time.Now().UnixNano())
}

func isValidInvoiceStatus(status string) bool {
	validStatuses := map[string]bool{
		"UNPAID":    true,
		"PAID":      true,
		"OVERDUE":   true,
		"CANCELLED": true,
	}

	return validStatuses[status]
}
