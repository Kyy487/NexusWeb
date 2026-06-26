package auth

import (
	"context"
	"errors"
	"time"

	"nexusweb-market/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type activityLogger interface {
	Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error
}

type service struct {
	repo   Repository
	cfg    *config.Config
	logger activityLogger
}

func NewService(repo Repository, cfg *config.Config, logger activityLogger) Service {
	return &service{repo: repo, cfg: cfg, logger: logger}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	existingUser, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	roleID, err := s.repo.FindRoleByName(ctx, "CUSTOMER")
	if err != nil {
		return nil, errors.New("customer role not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		RoleID:       roleID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.RoleName = "CUSTOMER"

	if s.logger != nil {
		_ = s.logger.Log(ctx, createdUser.ID, "AUTH", "REGISTER", "User registered successfully", "")
	}

	return &UserResponse{
		ID:     createdUser.ID,
		Name:   createdUser.Name,
		Email:  createdUser.Email,
		Phone:  createdUser.Phone,
		Role:   createdUser.RoleName,
		Status: createdUser.Status,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if user.Status != "ACTIVE" {
		return nil, errors.New("user account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	if s.logger != nil {
		_ = s.logger.Log(ctx, user.ID, "AUTH", "LOGIN", "User logged in successfully", "")
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Phone:  user.Phone,
			Role:   user.RoleName,
			Status: user.Status,
		},
	}, nil
}

func (s *service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.RoleName,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpiredHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.cfg.JWTSecret))
}
