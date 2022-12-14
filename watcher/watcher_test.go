package watcher_test

import (
	"fmt"
	"github.com/cpustejovsky/giles/watcher"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var configFilePath string
var badConfigFilePath string
var buildPath string
var binaryPath string

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	root = filepath.Dir(root)
	configFilePath = filepath.Join(root, "./test/config.yaml")
	badConfigFilePath = filepath.Join(root, "./test/badconfig.yaml")
	binaryPath = filepath.Join(root, "./test/test-binary")
	buildPath = filepath.Join(root, "build.sh")
}

func TestNewWatcher(t *testing.T) {
	t.Run("returns ar error if .yaml is not at end of file path", func(t *testing.T) {
		watcher, err := watcher.NewWatcher("foo/bar.json", buildPath)
		assert.Error(t, err)
		if watcher != nil {
			watcher.CloseWatcher()
		}
	})
	t.Run("returns error if file path not found", func(t *testing.T) {
		watcher, err := watcher.NewWatcher("foo/bar.yaml", buildPath)
		assert.Error(t, err)
		if watcher != nil {
			watcher.CloseWatcher()
		}
	})
	t.Run("returns no error if file path is found", func(t *testing.T) {
		watcher, err := watcher.NewWatcher(configFilePath, buildPath)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
	})
}

func TestWatcher_Close(t *testing.T) {
	t.Run("Watcher closes without error", func(t *testing.T) {
		watcher, err := watcher.NewWatcher(configFilePath, buildPath)
		assert.Nil(t, err)
		watcher.Start()
		err = watcher.CloseWatcher()
		assert.Nil(t, err)
	})
}

func TestWatcher_Start(t *testing.T) {
	t.Run("Start services without error", func(t *testing.T) {
		watcher, err := watcher.NewWatcher(configFilePath, buildPath)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
		watcher.Start()
		select {
		case err := <-watcher.ErrorChan:
			assert.Nil(t, err)
		case <-time.After(50 * time.Millisecond):
			return
		}
	})
	t.Run("Start services with error", func(t *testing.T) {
		watcher, err := watcher.NewWatcher(badConfigFilePath, buildPath)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
		watcher.Start()
		select {
		case err := <-watcher.ErrorChan:
			t.Log(err.Error())
			assert.Error(t, err)
			assert.NotNil(t, err)
		}
	})
}

func TestWatcher_Watch(t *testing.T) {
	t.Run("Watch files without error", func(t *testing.T) {
		watcher, err := watcher.NewWatcher(configFilePath, buildPath)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
		go watcher.Watch()
		select {
		case err := <-watcher.ErrorChan:
			assert.Nil(t, err)
		case <-time.After(50 * time.Millisecond):
			return
		}
	})
}

func TestWatcher_Build(t *testing.T) {
	name := "one"
	randomNum := "420"
	t.Run("Build passes with correct build path and arguments", func(t *testing.T) {

		wantBin := fmt.Sprintf("./tmp/builds/%s-%s", name, randomNum)
		bin, err := watcher.Build(buildPath, []string{"/home/cpustejovsky/go/src/giles/test/two", name, randomNum})
		assert.Nil(t, err)
		assert.Equal(t, wantBin, bin)
	})
	t.Run("Build fails with incorrect build path", func(t *testing.T) {
		bin, err := watcher.Build("", []string{"/home/cpustejovsky/go/src/giles/test/two", name, randomNum})
		assert.Error(t, err)
		assert.Equal(t, "", bin)
	})
	t.Run("Build fails with incorrect arguments", func(t *testing.T) {
		bin, err := watcher.Build(buildPath, []string{"/home/cputwo"})
		assert.Error(t, err)
		t.Log(err)
		assert.Equal(t, "", bin)
	})
}

func TestWatcher_Run(t *testing.T) {
	watcher, err := watcher.NewWatcher(configFilePath, buildPath)
	assert.Nil(t, err)
	defer watcher.CloseWatcher()
	var wg sync.WaitGroup
	wg.Add(1)
	err = watcher.Run(&wg, binaryPath)
	assert.Nil(t, err)
}
