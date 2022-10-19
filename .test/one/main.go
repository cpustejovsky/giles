package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "One",
	}

	log, _ := config.Build()
	sugar := log.Sugar()
	serverErrors := make(chan error, 1)
	go func() {
		sugar.Infow("startup", "status", "api router started")
		serverErrors <- http.ListenAndServe(":8001", nil)
	}()

	select {
	case err := <-serverErrors:
		sugar.Errorf("server error: %w", err)
	}
	sugar.Infow("stopping service")
}
