package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": "two",
	}

	l, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	log := l.Sugar()
	defer log.Sync()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("api router started using zap")
		serverErrors <- http.ListenAndServe(":8003", nil)
	}()

	select {
	case err := <-serverErrors:
		log.Infof("server error: %v\n", err)
		os.Exit(1)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
	}
}
