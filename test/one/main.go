package main

import (
	"log"
	"net/http"
)

func main() {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("api router started")
		serverErrors <- http.ListenAndServe(":8001", nil)
	}()

	select {
	case err := <-serverErrors:
		log.Printf("server error: %v\n", err)
	}
	log.Println("stopping service")
}
