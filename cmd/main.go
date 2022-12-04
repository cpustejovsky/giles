package main

import (
	"github.com/cpustejovsky/giles"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := os.Args[1]
	if config == "" {
		log.Println("please provide file path to yaml configuration file")
		os.Exit(1)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//pass in path to configuration yaml file
	w, err := giles.NewWatcher(config)
	if err != nil {
		log.Println("Error\t", err)
		os.Exit(1)
	}
	defer w.Close()

	//Start watching for file changes
	go w.Watch()
	//Start services
	w.Start()

	select {
	case err := <-w.ErrorChan:
		if err != nil {
			log.Println("Error\t", err)
			os.Exit(1)
		}
	case sig := <-sigs:
		log.Println("Signal\t", sig)
		os.Exit(0)
	}

}
