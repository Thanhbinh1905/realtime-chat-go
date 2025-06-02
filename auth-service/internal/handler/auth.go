package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/middleware"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/service"
	"go.uber.org/zap"

	"github.com/Thanhbinh1905/realtime-chat/shared/logger"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)

	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LogError("Failed to decode request body", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(r.Context(), &input); err != nil {
		logger.LogError("Failed to register user", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
	logger.LogInfo("User registration successful", zap.String("username", input.Username))
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LogError("Failed to decode request body", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(r.Context(), &input)
	if err != nil {
		logger.LogError("Failed to login user", err)
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	logger.LogInfo("User login successful", zap.String("username", input.Username))
}

func (h *authHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	curentUserID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || curentUserID == "" {
		logger.LogError("Unauthorized access", nil)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDToFetch := r.URL.Query().Get("id")
	if userIDToFetch == "" {
		userIDToFetch = curentUserID
	}

	user, err := h.service.GetUserByID(ctx, userIDToFetch)
	if err != nil {
		logger.LogError("Failed to get user by ID", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	response := model.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
