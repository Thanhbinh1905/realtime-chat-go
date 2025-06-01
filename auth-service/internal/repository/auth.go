package repository

import (
	"context"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, request *model.LoginRequest) (*model.User, error)

	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)",
		user.ID, user.Username, user.Email, user.Password,
	)

	return err
}

func (r *authRepository) Login(ctx context.Context, request *model.LoginRequest) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, "SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE username = $1", request.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, "SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE id = $1", userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
