package giles

import (
	"github.com/Docker/docker/pkg/filenotify"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
	pids        []int
	logger      zap.SugaredLogger
	rootPath    string
	buildPath   string
}

func New(rootPath string, l zap.SugaredLogger) (*watcher, error) {
	buildpath := filepath.Join(rootPath, "bin/tmp_build_go_file.sh")

	fileWatcher := filenotify.NewPollingWatcher()
	w := watcher{
		fileWatcher: fileWatcher,
		pids:        []int{},
		logger:      l,
		buildPath:   buildpath,
		rootPath:    rootPath,
	}
	return &w, nil
}

// Close wraps fileWatcher.Close
func (w *watcher) Close() error {
	return w.fileWatcher.Close()
}

/*
addPath is the WalkFunc used in AddPaths
It makes sure the path is a Service and then adds that Service to the embedded fileWatcher
*/
func (w *watcher) addPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		w.logger.Error(err)
		return err
	}
	if info.IsDir() {
		err = w.fileWatcher.Add(path)
		if err != nil {
			w.logger.Error(err)
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
			w.logger.Errorw("Error adding paths to watcher", "message", err, "path", path)
			return err
		}
	}
	return nil
}

// Watch watches for Events from the embedded fileWatcher and runs Restart
func (w *watcher) Watch(services []Service) {
	for {
		select {
		case _, ok := <-w.fileWatcher.Events():
			if !ok {
				return
			}
			w.Restart(services)
		case err, ok := <-w.fileWatcher.Errors():
			if !ok {
				return
			}
			w.logger.Error("error: ", err)
		}
	}
}

func (w *watcher) Restart(services []Service) {

}
