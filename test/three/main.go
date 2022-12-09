package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	serverErrors := make(chan error, 1)
	go func() {
		logrus.Info("api router started using logrus")
		serverErrors <- http.ListenAndServe(":8002", nil)
	}()

	select {
	case err := <-serverErrors:
		logrus.Error("server error: %v\n", err)
	}
	logrus.Info("stopping service")
}
