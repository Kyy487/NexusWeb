package dto

type OrderFileResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	UploadedBy string  `json:"uploaded_by"`
	FileName   string  `json:"file_name"`
	FileURL    string  `json:"file_url"`
	FileType   *string `json:"file_type"`
	FileSize   int64   `json:"file_size"`
	CreatedAt  string  `json:"created_at"`
}