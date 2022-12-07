package watcher

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Docker/docker/pkg/filenotify"
)

// Service provides us with a Name to use to identify the service when it's binary is built and run along with the path that it will run
type Service = struct {
	Name string
	Path string
}

// Watcher interface handles the only external facing filenotify.FileWatcher method that is being used outside this module
type Watcher interface {
	Close()
}

/*
watcher embeds filenotify.FileWatcher and contains:
a services slice to store services watcher will run
a pids slice for closing active processes
a sync.WaitGroup to determine when services have been started
an ErrorChan channel to pass errors to if goroutines fail
a buildPath variable to store location of the build bash script
*/
type watcher struct {
	filenotify.FileWatcher
	services  []Service
	buildPath string
	pids      []int
	wg        *sync.WaitGroup
	ErrorChan chan error
}

// NewWatcher takes a filePath pointing to the yaml config file and returns a watcher
func NewWatcher(filePath, buildpath string) (*watcher, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	wd = filepath.Dir(wd)
	log.Println("buildpath", buildpath)
	var wg sync.WaitGroup
	w := watcher{
		services:    []Service{},
		FileWatcher: filenotify.NewPollingWatcher(),
		pids:        []int{},
		buildPath:   buildpath,
		wg:          &wg,
		ErrorChan:   make(chan error, 3),
	}
	err = w.read(filePath)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// CloseWatcher removes tmp builds directory and closes embedded filewatcher
func (w *watcher) CloseWatcher() error {
	defer w.Close()
	defer w.stop()
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory:\t %v", err)
	}
	err = os.RemoveAll(filepath.Join(wd, "tmp"))
	if err != nil {
		return fmt.Errorf("error removing temp directory:\t %v", err)
	}
	if err != nil {
		return fmt.Errorf("error killing PIDS:\t %v", err)
	}
	return nil
}

/*
addPathWalkFunc is the WalkFunc used in AddPaths
It makes sure the path is a Service and then adds that Service to the embedded fileWatcher
*/
func (w *watcher) addPathWalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = w.Add(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPaths takes the paths provided and walks through them all to make sure all files contained within are watched
func (w *watcher) addPath(path string) error {
	return filepath.Walk(path, w.addPathWalkFunc)
}

// Watch watches for Events from the embedded fileWatcher and runs restart
func (w *watcher) Watch() {
	for {
		select {
		case _, ok := <-w.Events():
			//TODO: determine if these !ok returns should have specific errors
			if !ok {
				return
			}
			w.restart()
		case err, ok := <-w.Errors():
			if !ok {
				return
			}
			w.ErrorChan <- err
		}
	}
}

// restart runs stop then Start
func (w *watcher) restart() {
	err := w.stop()
	if err != nil {
		w.ErrorChan <- err
		return
	}
	w.Start()
}

// stop ranges over the pids that run adds to the watcher and kills each process. it then empties the pid slice
func (w *watcher) stop() error {
	for _, pid := range w.pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		err = proc.Kill()
		if err != nil {
			log.Printf("error killing process %v:\t%v\n", pid, err)
			return err
		}
		log.Printf("PID Killed:\t%v\n", pid)
	}
	w.pids = []int{}
	return nil
}

// Start ranges over a slice of Service and ranges over the Children of the Service, running BuildAndRun for each child in a goroutine
func (w *watcher) Start() {
	log.Println("Starting services")
	for _, service := range w.services {
		log.Println("Starting service: ", service.Name)
		w.wg.Add(1)
		randomNumber := rand.Intn(100)
		args := []string{service.Path, service.Name, strconv.Itoa(randomNumber)}
		go w.buildAndRun(w.buildPath, args)
	}
	log.Println("waiting")
	w.wg.Wait()
	log.Println("done waiting")
}

func (w *watcher) buildAndRun(buildpath string, args []string) {
	binary, err := build(buildpath, args)
	if err != nil {
		log.Println(err)
		w.ErrorChan <- err
	}
	err = w.run(binary)
	if err != nil {
		log.Println(err)
		w.ErrorChan <- err
	}
}

// build takes the buildpath to know what build script to run and any additional arguments to pass in
func build(buildpath string, args []string) (string, error) {
	output, err := exec.Command(buildpath, args...).Output()
	if err != nil {
		return "", errors.New(string(output))
	}
	binary := strings.TrimSpace(string(output))
	return binary, nil
}

// run creates a command from the binary path provided, sets up stdout and stderr, starts the command, and appends the process's Pid
func (w *watcher) run(binary string) error {
	defer w.wg.Done()
	cmd := exec.Command(binary)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	multi := io.MultiReader(stdout, stderr)
	in := bufio.NewScanner(multi)
	if err = cmd.Start(); err != nil {
		return err
	}
	w.pids = append(w.pids, cmd.Process.Pid)
	go func(input *bufio.Scanner) {
		for in.Scan() {
			log.Printf(input.Text()) // write each line to your log, or anything you need
		}
		if err = input.Err(); err != nil {
			w.ErrorChan <- err
		}
	}(in)
	return nil
}
