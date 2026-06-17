package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/category/dto"
	"nexusweb-market/backend/internal/modules/category/model"
	"nexusweb-market/backend/internal/modules/category/repository"

	"github.com/jackc/pgx/v5"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]dto.CategoryResponse, error)
	GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error)
	Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateCategoryRequest) error
	Delete(ctx context.Context, id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := []dto.CategoryResponse{}

	for _, category := range categories {
		responses = append(responses, mapCategoryToResponse(&category))
	}

	return responses, nil
}

func (s *categoryService) GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	response := mapCategoryToResponse(category)
	return &response, nil
}

func (s *categoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	category := &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      status,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	response := mapCategoryToResponse(category)
	return &response, nil
}

func (s *categoryService) Update(ctx context.Context, id string, req dto.UpdateCategoryRequest) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("category not found")
		}
		return err
	}

	category := &model.Category{
		ID:          id,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      req.Status,
	}

	return s.repo.Update(ctx, category)
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("category not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func mapCategoryToResponse(category *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Status:      category.Status,
	}
}