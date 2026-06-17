package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/orderrequirement/dto"
	"nexusweb-market/backend/internal/modules/orderrequirement/model"
	"nexusweb-market/backend/internal/modules/orderrequirement/repository"
)

type RequirementService interface {
	GetByOrderID(ctx context.Context, orderID string) ([]dto.RequirementResponse, error)
	GetByID(ctx context.Context, id string) (*dto.RequirementResponse, error)
	Create(ctx context.Context, orderID string, req dto.CreateRequirementRequest) (*dto.RequirementResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateRequirementRequest) (*dto.RequirementResponse, error)
	Delete(ctx context.Context, id string) error
}

type requirementService struct {
	repo repository.RequirementRepository
}

func NewRequirementService(repo repository.RequirementRepository) RequirementService {
	return &requirementService{repo: repo}
}

func (s *requirementService) GetByOrderID(ctx context.Context, orderID string) ([]dto.RequirementResponse, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	requirements, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RequirementResponse, 0, len(requirements))
	for _, requirement := range requirements {
		responses = append(responses, toRequirementResponse(requirement))
	}

	return responses, nil
}

func (s *requirementService) GetByID(ctx context.Context, id string) (*dto.RequirementResponse, error) {
	if id == "" {
		return nil, errors.New("requirement id is required")
	}

	requirement, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toRequirementResponse(*requirement)
	return &response, nil
}

func (s *requirementService) Create(ctx context.Context, orderID string, req dto.CreateRequirementRequest) (*dto.RequirementResponse, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	requirement := &model.Requirement{
		OrderID:  orderID,
		Question: req.Question,
		Answer:   req.Answer,
	}

	if err := s.repo.Create(ctx, requirement); err != nil {
		return nil, err
	}

	fullRequirement, err := s.repo.FindByID(ctx, requirement.ID)
	if err != nil {
		return nil, err
	}

	response := toRequirementResponse(*fullRequirement)
	return &response, nil
}

func (s *requirementService) Update(ctx context.Context, id string, req dto.UpdateRequirementRequest) (*dto.RequirementResponse, error) {
	if id == "" {
		return nil, errors.New("requirement id is required")
	}

	requirement := &model.Requirement{
		ID:       id,
		Question: req.Question,
		Answer:   req.Answer,
	}

	if err := s.repo.Update(ctx, requirement); err != nil {
		return nil, err
	}

	fullRequirement, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toRequirementResponse(*fullRequirement)
	return &response, nil
}

func (s *requirementService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("requirement id is required")
	}

	return s.repo.Delete(ctx, id)
}

func toRequirementResponse(requirement model.Requirement) dto.RequirementResponse {
	return dto.RequirementResponse{
		ID:       requirement.ID,
		OrderID:  requirement.OrderID,
		Question: requirement.Question,
		Answer:   requirement.Answer,
	}
}