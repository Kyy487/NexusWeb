package dto

type PaymentResponse struct {
	ID              string  `json:"id"`
	InvoiceID       string  `json:"invoice_id"`
	Amount          float64 `json:"amount"`
	PaymentMethod   *string `json:"payment_method"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentProofURL *string `json:"payment_proof_url"`
	PaidAt          *string `json:"paid_at"`
	VerifiedBy      *string `json:"verified_by"`
}

type CreatePaymentRequest struct {
	InvoiceID       string  `json:"invoice_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	PaymentMethod   string  `json:"payment_method"`
	PaymentProofURL string  `json:"payment_proof_url"`
}

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" binding:"required"`
	VerifiedBy    string `json:"verified_by"`
}
type WhatsAppPaymentResponse struct {
	InvoiceID   string  `json:"invoice_id"`
	Amount      float64 `json:"amount"`
	WhatsAppURL string  `json:"whatsapp_url"`
}