package giles_test

import (
	"github.com/cpustejovsky/giles"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	yamlLocation := filepath.Join(wd, "./test/test.yaml")
	t.Run("returns no error if .yaml at end of file path", func(t *testing.T) {
		watcher, err := giles.NewWatcher([]giles.Service{})
		assert.Nil(t, err)
		defer watcher.Close()
		err = watcher.Read(yamlLocation)
		assert.Nilf(t, err, "%v", err)
	})
	t.Run("returns ar error if .yaml is not at end of file path", func(t *testing.T) {
		watcher, err := giles.NewWatcher([]giles.Service{})
		assert.Nil(t, err)
		defer watcher.Close()
		err = watcher.Read("foo/bar.json")
		assert.Error(t, err)
	})
	t.Run("returns error if file path not found", func(t *testing.T) {
		watcher, err := giles.NewWatcher([]giles.Service{})
		assert.Nil(t, err)
		defer watcher.Close()
		err = watcher.Read("foo/bar.yaml")
		assert.Error(t, err)
	})
	t.Run("returns no error if file path is found", func(t *testing.T) {
		watcher, err := giles.NewWatcher([]giles.Service{})
		assert.Nil(t, err)
		defer watcher.Close()
		err = watcher.Read(yamlLocation)
		assert.Nil(t, err)
	})
}
