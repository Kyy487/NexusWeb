package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/orderprogress/dto"
	"nexusweb-market/backend/internal/modules/orderprogress/model"
	"nexusweb-market/backend/internal/modules/orderprogress/repository"
)

type ProgressService interface {
	GetByOrderID(ctx context.Context, orderID string) ([]dto.ProgressResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ProgressResponse, error)
	Create(ctx context.Context, orderID string, req dto.CreateProgressRequest) (*dto.ProgressResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateProgressRequest) (*dto.ProgressResponse, error)
	Delete(ctx context.Context, id string) error
}

type progressService struct {
	repo repository.ProgressRepository
}

func NewProgressService(repo repository.ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) GetByOrderID(ctx context.Context, orderID string) ([]dto.ProgressResponse, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	progressList, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ProgressResponse, 0, len(progressList))
	for _, progress := range progressList {
		responses = append(responses, toProgressResponse(progress))
	}

	return responses, nil
}

func (s *progressService) GetByID(ctx context.Context, id string) (*dto.ProgressResponse, error) {
	if id == "" {
		return nil, errors.New("progress id is required")
	}

	progress, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toProgressResponse(*progress)
	return &response, nil
}

func (s *progressService) Create(ctx context.Context, orderID string, req dto.CreateProgressRequest) (*dto.ProgressResponse, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	if req.ProgressPercentage < 0 || req.ProgressPercentage > 100 {
		return nil, errors.New("progress percentage must be between 0 and 100")
	}

	progress := &model.Progress{
		OrderID:            orderID,
		Title:              req.Title,
		Description:        req.Description,
		ProgressPercentage: req.ProgressPercentage,
		CreatedBy:          req.CreatedBy,
	}

	if err := s.repo.Create(ctx, progress); err != nil {
		return nil, err
	}

	fullProgress, err := s.repo.FindByID(ctx, progress.ID)
	if err != nil {
		return nil, err
	}

	response := toProgressResponse(*fullProgress)
	return &response, nil
}

func (s *progressService) Update(ctx context.Context, id string, req dto.UpdateProgressRequest) (*dto.ProgressResponse, error) {
	if id == "" {
		return nil, errors.New("progress id is required")
	}

	if req.ProgressPercentage < 0 || req.ProgressPercentage > 100 {
		return nil, errors.New("progress percentage must be between 0 and 100")
	}

	progress := &model.Progress{
		ID:                 id,
		Title:              req.Title,
		Description:        req.Description,
		ProgressPercentage: req.ProgressPercentage,
	}

	if err := s.repo.Update(ctx, progress); err != nil {
		return nil, err
	}

	fullProgress, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toProgressResponse(*fullProgress)
	return &response, nil
}

func (s *progressService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("progress id is required")
	}

	return s.repo.Delete(ctx, id)
}

func toProgressResponse(progress model.Progress) dto.ProgressResponse {
	return dto.ProgressResponse{
		ID:                 progress.ID,
		OrderID:            progress.OrderID,
		Title:              progress.Title,
		Description:        progress.Description,
		ProgressPercentage: progress.ProgressPercentage,
		CreatedBy:          progress.CreatedBy,
	}
}