package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/service"
	"go.uber.org/zap"

	"github.com/Thanhbinh1905/realtime-chat/shared/logger"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
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
