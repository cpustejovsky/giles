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

var rootPath string

type testService struct {
	name string
	port string
	path string
}

var testServices = []testService{
	{
		name: "one",
		port: "8001",
	},
	{
		name: "two",
		port: "8002",
	},
	{
		name: "three",
		port: "8003",
	},
}

func cleanUpTests() {
	err := os.RemoveAll(filepath.Join(rootPath, "tmp/builds"))
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootPath = root
}

func TestWatcher_AddPaths(t *testing.T) {

	watcher, err := giles.NewWatcher(rootPath)
	assert.Nil(t, err)
	defer watcher.Close()
	t.Run("Add existing directories to watch", func(t *testing.T) {
		err := watcher.AddPaths([]string{filepath.Join(rootPath, ".test")})
		assert.Nil(t, err)
	})
	t.Run("Add unknown directories to watch", func(t *testing.T) {
		err := watcher.AddPaths([]string{filepath.Join(rootPath, "foobarbaz")})
		assert.Error(t, err)
	})
	t.Cleanup(cleanUpTests)
}

func TestWatcher_Start(t *testing.T) {
	t.Run("Run services without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(rootPath)
		assert.Nil(t, err)
		defer watcher.Close()
		var services []giles.Service
		for _, service := range testServices {
			services = append(services, giles.Service{
				Name: service.name,
				Path: service.path,
			})
		}
		watcher.Start(services)
		select {
		case err := <-watcher.ErrorChan:
			assert.Nil(t, err)
		case <-time.After(50 * time.Millisecond):
			return
		}
	})
	t.Run("Run services with error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(rootPath)
		assert.Nil(t, err)
		defer watcher.Close()
		watcher.Start([]giles.Service{{Name: "foobar", Path: "foobarbaz"}})
		select {
		case err := <-watcher.ErrorChan:
			assert.Error(t, err)
		}
	})
	t.Cleanup(cleanUpTests)
}
