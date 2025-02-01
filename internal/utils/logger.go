package utils

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	logger, err := zap.NewProduction() 
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	Log = logger
}
