package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/package/dto"
	"nexusweb-market/backend/internal/modules/package/model"
	"nexusweb-market/backend/internal/modules/package/repository"
)

type PackageService interface {
	GetAll(ctx context.Context) ([]dto.PackageResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PackageResponse, error)
	Create(ctx context.Context, req dto.CreatePackageRequest) (*dto.PackageResponse, error)
	Update(ctx context.Context, id string, req dto.UpdatePackageRequest) (*dto.PackageResponse, error)
	Delete(ctx context.Context, id string) error
}

type packageService struct {
	repo repository.PackageRepository
}

func NewPackageService(repo repository.PackageRepository) PackageService {
	return &packageService{repo: repo}
}

func (s *packageService) GetAll(ctx context.Context) ([]dto.PackageResponse, error) {
	packages, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.PackageResponse
	for _, pkg := range packages {
		responses = append(responses, toPackageResponse(pkg))
	}

	return responses, nil
}

func (s *packageService) GetByID(ctx context.Context, id string) (*dto.PackageResponse, error) {
	if id == "" {
		return nil, errors.New("package id is required")
	}

	pkg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toPackageResponse(*pkg)
	return &response, nil
}

func (s *packageService) Create(ctx context.Context, req dto.CreatePackageRequest) (*dto.PackageResponse, error) {
	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	pkg := &model.Package{
		ServiceID:      req.ServiceID,
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		RevisionCount:  req.RevisionCount,
		DeliveryDays:   req.DeliveryDays,
		Features:       req.Features,
		Status:         req.Status,
	}

	err := s.repo.Create(ctx, pkg)
	if err != nil {
		return nil, err
	}

	fullPkg, err := s.repo.FindByID(ctx, pkg.ID)
	if err != nil {
		return nil, err
	}

	response := toPackageResponse(*fullPkg)
	return &response, nil
}

func (s *packageService) Update(ctx context.Context, id string, req dto.UpdatePackageRequest) (*dto.PackageResponse, error) {
	if id == "" {
		return nil, errors.New("package id is required")
	}

	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	pkg := &model.Package{
		ID:             id,
		ServiceID:      req.ServiceID,
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		RevisionCount:  req.RevisionCount,
		DeliveryDays:   req.DeliveryDays,
		Features:       req.Features,
		Status:         req.Status,
	}

	err := s.repo.Update(ctx, pkg)
	if err != nil {
		return nil, err
	}

	fullPkg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := toPackageResponse(*fullPkg)
	return &response, nil
}

func (s *packageService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("package id is required")
	}

	return s.repo.Delete(ctx, id)
}

func toPackageResponse(pkg model.Package) dto.PackageResponse {
	return dto.PackageResponse{
		ID:             pkg.ID,
		ServiceID:      pkg.ServiceID,
		ServiceName:    pkg.ServiceName,
		Name:           pkg.Name,
		Description:    pkg.Description,
		Price:          pkg.Price,
		RevisionCount:  pkg.RevisionCount,
		DeliveryDays:   pkg.DeliveryDays,
		Features:       pkg.Features,
		Status:         pkg.Status,
	}
}