package main

import (
	"go.uber.org/zap"
)

func main() {
	// Khởi tạo logger cho môi trường phát triển
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // Đảm bảo tất cả log được ghi ra trước khi thoát

	logger.Info("Đây là một log info từ logger phát triển",
		zap.String("key", "value"),
		zap.Int("count", 123),
	)
	logger.Debug("Đây là một log debug", zap.Bool("enabled", true))
	logger.Warn("Đây là một cảnh báo", zap.Error(nil))
	logger.Error("Đây là một lỗi", zap.Error(nil)) // Sử dụng nil cho ví dụ
}
