package service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"nexusweb-market/backend/internal/modules/orderfile/dto"
	"nexusweb-market/backend/internal/modules/orderfile/model"
	"nexusweb-market/backend/internal/modules/orderfile/repository"
)

type OrderFileService interface {
	SaveFile(ctx context.Context, orderID string, uploadedBy string, fileType string, fileHeader *multipart.FileHeader, fileURL string) (*dto.OrderFileResponse, error)
	GetByOrderID(ctx context.Context, orderID string) ([]dto.OrderFileResponse, error)
	Delete(ctx context.Context, id string) error
}

type orderFileService struct {
	repo repository.OrderFileRepository
}

func NewOrderFileService(repo repository.OrderFileRepository) OrderFileService {
	return &orderFileService{repo: repo}
}

func (s *orderFileService) SaveFile(ctx context.Context, orderID string, uploadedBy string, fileType string, fileHeader *multipart.FileHeader, fileURL string) (*dto.OrderFileResponse, error) {
	if orderID == "" {
		return nil, errors.New("order_id is required")
	}

	if uploadedBy == "" {
		return nil, errors.New("uploaded_by is required")
	}

	if fileHeader == nil {
		return nil, errors.New("file is required")
	}

	if !isAllowedFile(fileHeader.Filename) {
		return nil, errors.New("file type not allowed")
	}

	var fileTypePtr *string
	if fileType != "" {
		fileTypePtr = &fileType
	}

	orderFile := &model.OrderFile{
		OrderID:    orderID,
		UploadedBy: uploadedBy,
		FileName:   fileHeader.Filename,
		FileURL:    fileURL,
		FileType:   fileTypePtr,
		FileSize:   fileHeader.Size,
	}

	if err := s.repo.Create(ctx, orderFile); err != nil {
		return nil, err
	}

	response := toOrderFileResponse(*orderFile)
	return &response, nil
}

func (s *orderFileService) GetByOrderID(ctx context.Context, orderID string) ([]dto.OrderFileResponse, error) {
	if orderID == "" {
		return nil, errors.New("order_id is required")
	}

	files, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderFileResponse, 0, len(files))
	for _, file := range files {
		responses = append(responses, toOrderFileResponse(file))
	}

	return responses, nil
}

func (s *orderFileService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("file id is required")
	}

	return s.repo.Delete(ctx, id)
}

func toOrderFileResponse(file model.OrderFile) dto.OrderFileResponse {
	return dto.OrderFileResponse{
		ID:         file.ID,
		OrderID:    file.OrderID,
		UploadedBy: file.UploadedBy,
		FileName:   file.FileName,
		FileURL:    file.FileURL,
		FileType:   file.FileType,
		FileSize:   file.FileSize,
		CreatedAt:  file.CreatedAt.Format(time.RFC3339),
	}
}

func isAllowedFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".zip":  true,
		".rar":  true,
	}

	return allowed[ext]
}