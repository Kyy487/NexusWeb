package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/service/dto"
	"nexusweb-market/backend/internal/modules/service/model"
	"nexusweb-market/backend/internal/modules/service/repository"

	"github.com/jackc/pgx/v5"
)

type ServiceService interface {
	GetAll(ctx context.Context) ([]dto.ServiceResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ServiceResponse, error)
	Create(ctx context.Context, req dto.CreateServiceRequest) (*dto.ServiceResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateServiceRequest) error
	Delete(ctx context.Context, id string) error
}

type serviceService struct {
	repo repository.ServiceRepository
}

func NewServiceService(repo repository.ServiceRepository) ServiceService {
	return &serviceService{repo: repo}
}

func (s *serviceService) GetAll(ctx context.Context) ([]dto.ServiceResponse, error) {
	services, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := []dto.ServiceResponse{}

	for _, service := range services {
		responses = append(responses, mapServiceToResponse(&service))
	}

	return responses, nil
}

func (s *serviceService) GetByID(ctx context.Context, id string) (*dto.ServiceResponse, error) {
	service, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}

	response := mapServiceToResponse(service)
	return &response, nil
}

func (s *serviceService) Create(ctx context.Context, req dto.CreateServiceRequest) (*dto.ServiceResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	estimatedDays := req.EstimatedDays
	if estimatedDays <= 0 {
		estimatedDays = 1
	}

	service := &model.Service{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		BasePrice:     req.BasePrice,
		EstimatedDays: estimatedDays,
		Status:        status,
	}

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, err
	}

	createdService, err := s.repo.FindByID(ctx, service.ID)
	if err != nil {
		return nil, err
	}

	response := mapServiceToResponse(createdService)
	return &response, nil
}

func (s *serviceService) Update(ctx context.Context, id string, req dto.UpdateServiceRequest) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("service not found")
		}
		return err
	}

	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	estimatedDays := req.EstimatedDays
	if estimatedDays <= 0 {
		estimatedDays = 1
	}

	service := &model.Service{
		ID:            id,
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		BasePrice:     req.BasePrice,
		EstimatedDays: estimatedDays,
		Status:        status,
	}

	return s.repo.Update(ctx, service)
}

func (s *serviceService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("service not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func mapServiceToResponse(service *model.Service) dto.ServiceResponse {
	return dto.ServiceResponse{
		ID:            service.ID,
		CategoryID:    service.CategoryID,
		CategoryName:  service.CategoryName,
		Name:          service.Name,
		Slug:          service.Slug,
		Description:   service.Description,
		BasePrice:     service.BasePrice,
		EstimatedDays: service.EstimatedDays,
		Status:        service.Status,
	}
}