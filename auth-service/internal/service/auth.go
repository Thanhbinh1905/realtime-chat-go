package service

import (
	"context"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/repository"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/utils"

	"github.com/Thanhbinh1905/realtime-chat/shared/logger"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type AuthService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, req *model.LoginRequest) (model.LoginResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, user *model.User) error {
	if err := validate.Struct(user); err != nil {
		logger.LogError("Validation failed", err)
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError("Failed to hash password", err)
		return err
	}

	user.ID = model.GenerateUUID()
	user.Password = string(hashedPassword)
	user.CreatedAt = model.GetCurrentTime()
	user.UpdatedAt = user.CreatedAt

	return s.repo.Register(ctx, user)
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (model.LoginResponse, error) {
	if err := validate.Struct(req); err != nil {
		return model.LoginResponse{}, err
	}

	user, err := s.repo.Login(ctx, req)
	if err != nil {
		logger.LogError("Login failed", err)
		return model.LoginResponse{}, err
	}

	if user == nil {
		logger.LogError("User not found", nil)
		return model.LoginResponse{}, utils.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.LogError("Password mismatch", err)
		return model.LoginResponse{}, utils.ErrInvalidCredentials
	}

	idToken, err_id := utils.GenerateIDToken(user.ID.String())
	accessToken, err_access := utils.GenerateAccessToken(user.ID.String(), string(user.Role))
	refreshToken, err_refresh := utils.GenerateRefreshToken(user.ID.String())
	if err_id != nil || err_access != nil || err_refresh != nil {
		logger.LogError("Failed to generate JWT", err)
		return model.LoginResponse{}, err
	}
	return model.LoginResponse{
		IDToken:      idToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	if userId == "" {
		logger.LogError("User ID is empty", nil)
		return nil, utils.ErrInvalidUserID
	}

	user, err := s.repo.GetUserByID(ctx, userId)
	if err != nil {
		logger.LogError("Failed to get user by ID", err)
		return nil, err
	}

	if user == nil {
		logger.LogError("User not found", nil)
		return nil, utils.ErrUserNotFound
	}

	return user, nil
}
