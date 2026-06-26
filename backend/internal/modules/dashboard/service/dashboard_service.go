package service

import (
	"context"

	"nexusweb-market/backend/internal/modules/dashboard/dto"
	"nexusweb-market/backend/internal/modules/dashboard/repository"
)

type DashboardService interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	return s.repo.GetStats(ctx)
}