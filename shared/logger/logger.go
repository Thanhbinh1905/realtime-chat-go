package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(isDev bool) {
	var err error
	if isDev {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}
	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
}

func LogError(msg string, err error) {
	if err != nil {
		Log.Error(msg, zap.Error(err))
	}
}
func LogInfo(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}
func LogDebug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}
