package main

import (
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/config"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/db"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/handler"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/middleware"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/repository"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/service"
	"github.com/gorilla/mux"

	"net/http"

	"github.com/Thanhbinh1905/realtime-chat/shared/logger"
	"go.uber.org/zap"
)

func main() {
	// Khởi tạo logger cho môi trường phát triển
	logger.Init(true)

	// Tải cấu hình từ biến môi trường
	cfg := config.LoadConfig()

	// Kết nối database
	err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db.Pool)
	service := service.NewAuthService(repo)
	handler := handler.NewAuthHandler(service)

	mw := middleware.New(cfg)

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AuthService OK"))
	}).Methods("GET")

	api.HandleFunc("/register", handler.Register).Methods("POST")

	api.HandleFunc("/login", handler.Login).Methods("POST")

	api.Handle("/user", mw.JWTAuthMiddleware(http.HandlerFunc(handler.GetUserByID))).Methods("GET")

	api.Use(
		middleware.LoggingMiddleware,  // Middleware để ghi log
		middleware.RecoveryMiddleware, // Middleware để phục hồi từ panic
	)
	// Bắt đầu lắng nghe và phục vụ HTTP
	logger.Log.Info("Starting server", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}

	logger.Log.Info("Server started successfully", zap.String("port", cfg.Port))
	logger.Log.Info("Auth service is running", zap.String("port", cfg.Port))
	logger.Log.Info("Listening on port", zap.String("port", cfg.Port))

}
