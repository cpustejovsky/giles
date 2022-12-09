package main

import (
	"fmt"
	"net/http"
)

func main() {
	serverErrors := make(chan error, 1)
	go func() {
		fmt.Println("api router started using fmt")
		serverErrors <- http.ListenAndServe(":8001", nil)
	}()

	select {
	case err := <-serverErrors:
		fmt.Printf("server error: %v\n", err)
	}
	fmt.Println("stopping service")
}
