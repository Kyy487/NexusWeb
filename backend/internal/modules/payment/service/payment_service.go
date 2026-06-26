package service

import (
	"context"
	"errors"
	"time"
	"fmt"
	"net/url"
	"os"

	"nexusweb-market/backend/internal/modules/payment/dto"
	"nexusweb-market/backend/internal/modules/payment/model"
	"nexusweb-market/backend/internal/modules/payment/repository"
)

type PaymentService interface {
	GetAll(ctx context.Context) ([]dto.PaymentResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error)
	GetByInvoiceID(ctx context.Context, invoiceID string) ([]dto.PaymentResponse, error)
	Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	UpdateStatus(ctx context.Context, id string, req dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error)
	GetWhatsAppLink(ctx context.Context, paymentID string) (*dto.WhatsAppPaymentResponse, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) GetAll(ctx context.Context) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, toPaymentResponse(payment))
	}

	return responses, nil
}

func (s *paymentService) GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error) {
	if id == "" {
		return nil, errors.New("payment id is required")
	}

	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toPaymentResponse(*payment)
	return &response, nil
}

func (s *paymentService) GetByInvoiceID(ctx context.Context, invoiceID string) ([]dto.PaymentResponse, error) {
	if invoiceID == "" {
		return nil, errors.New("invoice id is required")
	}

	payments, err := s.repo.FindByInvoiceID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		responses = append(responses, toPaymentResponse(payment))
	}

	return responses, nil
}

func (s *paymentService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	var paymentMethod *string
	if req.PaymentMethod != "" {
		paymentMethod = &req.PaymentMethod
	}

	var paymentProofURL *string
	if req.PaymentProofURL != "" {
		paymentProofURL = &req.PaymentProofURL
	}

	payment := &model.Payment{
		InvoiceID:       req.InvoiceID,
		Amount:          req.Amount,
		PaymentMethod:   paymentMethod,
		PaymentStatus:   "PENDING",
		PaymentProofURL: paymentProofURL,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	fullPayment, err := s.repo.FindByID(ctx, payment.ID)
	if err != nil {
		return nil, err
	}

	response := toPaymentResponse(*fullPayment)
	return &response, nil
}

func (s *paymentService) UpdateStatus(ctx context.Context, id string, req dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error) {
	if id == "" {
		return nil, errors.New("payment id is required")
	}

	if !isValidPaymentStatus(req.PaymentStatus) {
		return nil, errors.New("invalid payment status")
	}

	if req.PaymentStatus == "PAID" && req.VerifiedBy == "" {
		return nil, errors.New("verified_by is required when payment status is PAID")
	}

	var verifiedBy *string
	if req.VerifiedBy != "" {
		verifiedBy = &req.VerifiedBy
	}

	if err := s.repo.UpdateStatus(ctx, id, req.PaymentStatus, verifiedBy); err != nil {
		return nil, err
	}

	if req.PaymentStatus == "PAID" {
		if err := s.repo.UpdateInvoiceAndOrderAfterPayment(ctx, id); err != nil {
			return nil, err
		}
	}

	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toPaymentResponse(*payment)
	return &response, nil
}

func toPaymentResponse(payment model.Payment) dto.PaymentResponse {
	var paidAt *string
	if payment.PaidAt != nil {
		formatted := payment.PaidAt.Format(time.RFC3339)
		paidAt = &formatted
	}

	return dto.PaymentResponse{
		ID:              payment.ID,
		InvoiceID:       payment.InvoiceID,
		Amount:          payment.Amount,
		PaymentMethod:   payment.PaymentMethod,
		PaymentStatus:   payment.PaymentStatus,
		PaymentProofURL: payment.PaymentProofURL,
		PaidAt:          paidAt,
		VerifiedBy:      payment.VerifiedBy,
	}
}

func isValidPaymentStatus(status string) bool {
	validStatuses := map[string]bool{
		"PENDING":  true,
		"PAID":     true,
		"FAILED":   true,
		"EXPIRED":  true,
		"REFUNDED": true,
	}

	return validStatuses[status]
}
func (s *paymentService) GetWhatsAppLink(ctx context.Context, paymentID string) (*dto.WhatsAppPaymentResponse, error) {
	if paymentID == "" {
		return nil, errors.New("payment id is required")
	}

	adminWhatsApp := os.Getenv("ADMIN_WHATSAPP")
	if adminWhatsApp == "" {
		adminWhatsApp = "6281234567890"
	}

	payment, err := s.repo.GetWhatsAppData(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf(
		`Halo Admin NexusWeb,

Saya ingin melakukan konfirmasi order dan pembayaran.

━━━━━━━━━━━━━━━━━━━━
DETAIL ORDER
━━━━━━━━━━━━━━━━━━━━

Invoice :
%s

Service :
%s

Package :
%s

Total Pembayaran :
Rp%.0f

━━━━━━━━━━━━━━━━━━━━
TINDAKAN SELANJUTNYA
━━━━━━━━━━━━━━━━━━━━

Mohon kirimkan informasi:

• Rekening Pembayaran
atau
• QRIS Pembayaran

Setelah pembayaran berhasil dilakukan, saya akan mengirimkan:

• Bukti Pembayaran
• Requirement Project
• File Pendukung (Logo, Dokumen, Referensi, dll)

Terima kasih.

Salam,
Customer NexusWeb`,
		payment.InvoiceID,
		"Konfirmasi Layanan NexusWeb",
		"Package sesuai invoice",
		payment.Amount,
	)

	whatsAppURL := fmt.Sprintf(
		"https://wa.me/%s?text=%s",
		adminWhatsApp,
		url.QueryEscape(message),
	)

	return &dto.WhatsAppPaymentResponse{
		InvoiceID:   payment.InvoiceID,
		Amount:      payment.Amount,
		WhatsAppURL: whatsAppURL,
	}, nil
}