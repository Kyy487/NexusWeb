package service

import (
	"context"
	"errors"

	"nexusweb-market/backend/internal/modules/user/dto"
	"nexusweb-market/backend/internal/modules/user/model"
	"nexusweb-market/backend/internal/modules/user/repository"

	"github.com/jackc/pgx/v5"
)

type UserService interface {
	GetProfile(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	UpdateUserStatus(ctx context.Context, id string, req dto.UpdateUserStatusRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := mapUserToResponse(user)
	return &response, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := []dto.UserResponse{}

	for _, user := range users {
		response := mapUserToResponse(&user)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *userService) UpdateUserStatus(ctx context.Context, id string, req dto.UpdateUserStatusRequest) error {
	if req.Status != "ACTIVE" && req.Status != "INACTIVE" && req.Status != "SUSPENDED" {
		return errors.New("invalid user status")
	}

	return s.repo.UpdateStatus(ctx, id, req.Status)
}

func mapUserToResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Phone:  user.Phone,
		Role:   user.RoleName,
		Status: user.Status,
	}
}