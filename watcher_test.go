package giles_test

import (
	"github.com/cpustejovsky/giles"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var configFilePath string
var badConfigFilePath string

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configFilePath = filepath.Join(root, "./test/config.yaml")
	badConfigFilePath = filepath.Join(root, "./test/badconfig.yaml")
}

func TestNewWatcher(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	yamlLocation := filepath.Join(wd, "./test/test.yaml")
	t.Run("returns ar error if .yaml is not at end of file path", func(t *testing.T) {
		watcher, err := giles.NewWatcher("foo/bar.json")
		assert.Error(t, err)
		if watcher != nil {
			watcher.CloseWatcher()
		}
	})
	t.Run("returns error if file path not found", func(t *testing.T) {
		watcher, err := giles.NewWatcher("foo/bar.yaml")
		assert.Error(t, err)
		if watcher != nil {
			watcher.CloseWatcher()
		}
	})
	t.Run("returns no error if file path is found", func(t *testing.T) {
		watcher, err := giles.NewWatcher(yamlLocation)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
	})
}

func TestWatcher_Close(t *testing.T) {
	t.Run("Watcher closes without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(configFilePath)
		assert.Nil(t, err)
		watcher.Start()
		err = watcher.CloseWatcher()
		assert.Nil(t, err)
	})
}

func TestWatcher_Start(t *testing.T) {
	t.Run("Start services without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(configFilePath)
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
		watcher, err := giles.NewWatcher(badConfigFilePath)
		assert.Nil(t, err)
		defer watcher.CloseWatcher()
		watcher.Start()
		select {
		case err := <-watcher.ErrorChan:
			assert.Error(t, err)
		}
	})
}

func TestWatcher_Watch(t *testing.T) {
	t.Run("Watch files without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(configFilePath)
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
