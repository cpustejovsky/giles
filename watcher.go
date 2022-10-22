package giles

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
a pids slice for closing active processes
a logger of type zap.SugaredLogger
a rootPath to determine the tmp Service to build to
*/
type watcher struct {
	fileWatcher filenotify.FileWatcher
	services    []Service
	pids        []int
	buildPath   string
	ErrorChan   chan error
}

func NewWatcher(services []Service) (*watcher, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	buildpath := filepath.Join(wd, "./build.sh")
	fileWatcher := filenotify.NewPollingWatcher()
	w := watcher{
		services:    services,
		fileWatcher: fileWatcher,
		pids:        []int{},
		buildPath:   buildpath,
		ErrorChan:   make(chan error),
	}
	return &w, nil
}

// Close removes tmp builds directory and closes embedded filewatcher
func (w *watcher) Close() error {
	wd, err := os.Getwd()
	err = os.RemoveAll(filepath.Join(wd, "tmp"))
	if err != nil {
		return err
	}
	err = w.stop()
	if err != nil {
		return err
	}
	return w.fileWatcher.Close()
}

/*
addPath is the WalkFunc used in AddPaths
It makes sure the path is a Service and then adds that Service to the embedded fileWatcher
*/
func (w *watcher) addPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = w.fileWatcher.Add(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPaths takes the paths provided and walks through them all to make sure all files contained within are watched
func (w *watcher) AddPaths(paths []string) error {
	for _, path := range paths {
		err := filepath.Walk(path, w.addPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Watch watches for Events from the embedded fileWatcher and runs restart
func (w *watcher) Watch() {
	for {
		select {
		case _, ok := <-w.fileWatcher.Events():
			//TODO: determine if these !ok returns should have specific errors
			if !ok {
				return
			}
			w.restart()
		case err, ok := <-w.fileWatcher.Errors():
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
			return err
		}
		log.Printf("PID Killed:\t%v\n", pid)
	}
	w.pids = []int{}
	return nil
}

// Start ranges over a slice of Service and ranges over the Children of the Service, running BuildAndRun for each child in a goroutine
func (w *watcher) Start() {
	for _, service := range w.services {
		path := service.Path
		name := service.Name
		randomNumber := rand.Intn(100)
		args := []string{path, name, strconv.Itoa(randomNumber)}
		go func() {
			binary, err := w.build(w.buildPath, args)
			if err != nil {
				w.ErrorChan <- err
			}
			w.ErrorChan <- w.run(binary)
		}()
	}
}

// build takes the buildpath to know what build script to run and any additional arguments to pass in
func (w *watcher) build(buildpath string, args []string) (string, error) {
	output, err := exec.Command(buildpath, args...).Output()
	if err != nil {
		return "", err
	}
	binary := strings.TrimSpace(string(output))
	return binary, nil
}

// run creates a command from the binary path provided, sets up stdout and stderr, starts the command, and appends the process's Pid
func (w *watcher) run(binary string) error {
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
	if err := cmd.Start(); err != nil {
		return err
	}
	w.pids = append(w.pids, cmd.Process.Pid)
	for in.Scan() {
		log.Printf(in.Text()) // write each line to your log, or anything you need
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	return nil
}
