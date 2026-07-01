package service

import (
	"context"
	"errors"
	"time"

	"nexusweb-market/backend/internal/modules/activitylog/dto"
	"nexusweb-market/backend/internal/modules/activitylog/model"
	"nexusweb-market/backend/internal/modules/activitylog/repository"
)

type ActivityLogService interface {
	Create(ctx context.Context, req dto.CreateActivityLogRequest) (*dto.ActivityLogResponse, error)
	GetAll(ctx context.Context, userID string, role string) ([]dto.ActivityLogResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]dto.ActivityLogResponse, error)
	Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error
}

type activityLogService struct {
	repo repository.ActivityLogRepository
}

func NewActivityLogService(repo repository.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) Create(ctx context.Context, req dto.CreateActivityLogRequest) (*dto.ActivityLogResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	if req.Module == "" {
		return nil, errors.New("module is required")
	}

	if req.Action == "" {
		return nil, errors.New("action is required")
	}

	var ipAddress *string
	if req.IPAddress != "" {
		ipAddress = &req.IPAddress
	}

	log := &model.ActivityLog{
		UserID:      req.UserID,
		Module:      req.Module,
		Action:      req.Action,
		Description: req.Description,
		IPAddress:   ipAddress,
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, err
	}

	response := toActivityLogResponse(*log)
	return &response, nil
}

func (s *activityLogService) GetAll(ctx context.Context, userID string, role string) ([]dto.ActivityLogResponse, error) {
	var logs []model.ActivityLog
	var err error

	if role == "CUSTOMER" {
		logs, err = s.repo.FindByUserID(ctx, userID)
	} else {
		logs, err = s.repo.FindAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	responses := make([]dto.ActivityLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, toActivityLogResponse(log))
	}

	return responses, nil
}

func (s *activityLogService) GetByUserID(ctx context.Context, userID string) ([]dto.ActivityLogResponse, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	logs, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ActivityLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, toActivityLogResponse(log))
	}

	return responses, nil
}

func (s *activityLogService) Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error {
	if userID == "" {
		return nil
	}

	if module == "" || action == "" {
		return nil
	}

	var ip *string
	if ipAddress != "" {
		ip = &ipAddress
	}

	log := &model.ActivityLog{
		UserID:      userID,
		Module:      module,
		Action:      action,
		Description: description,
		IPAddress:   ip,
	}

	return s.repo.Create(ctx, log)
}

func toActivityLogResponse(log model.ActivityLog) dto.ActivityLogResponse {
	return dto.ActivityLogResponse{
		ID:          log.ID,
		UserID:      log.UserID,
		Module:      log.Module,
		Action:      log.Action,
		Description: log.Description,
		IPAddress:   log.IPAddress,
		CreatedAt:   log.CreatedAt.Format(time.RFC3339),
	}
}