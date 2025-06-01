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
