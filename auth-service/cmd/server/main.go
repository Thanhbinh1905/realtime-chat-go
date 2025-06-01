package main

import (
	"net/http"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/config"
	"github.com/Thanhbinh1905/realtime-chat/shared/logger"
	"go.uber.org/zap"
)

func main() {
	// Khởi tạo logger cho môi trường phát triển
	logger.Init(true)

	// Tải cấu hình từ biến môi trường
	cfg := config.LoadConfig()

	logger.Log.Info("Starting Auth Service", zap.String("port", cfg.Port))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Auth Service is running"))
	})

	// Thêm đoạn này để chạy server, log lỗi nếu có
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
