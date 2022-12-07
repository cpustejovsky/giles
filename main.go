package main

import (
	"github.com/cpustejovsky/giles/watcher"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("please provide file path to yaml configuration file")
		os.Exit(1)
	}
	config := os.Args[1]
	wd, err := os.Getwd()
	buildPath := filepath.Join(wd, "build.sh")
	if err != nil {
		log.Println("Error\t", err)
		os.Exit(1)
	}
	log.Println(wd)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//pass in path to configuration yaml file
	w, err := watcher.NewWatcher(config, buildPath)
	if err != nil {
		log.Println("Error\t", err)
		os.Exit(1)
	}
	defer w.Close()

	//Start services
	w.Start()
	//Start watching for file changes
	go w.Watch()

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
