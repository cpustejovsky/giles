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
	var services []giles.Service
	for _, service := range testServices {
		services = append(services, giles.Service{
			Name: service.name,
			Path: service.path,
		})
	}
	watcher, err := giles.NewWatcher(services)
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
	var services []giles.Service
	for _, service := range testServices {
		services = append(services, giles.Service{
			Name: service.name,
			Path: service.path,
		})
	}
	t.Run("Start services without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(services)
		assert.Nil(t, err)
		defer watcher.Close()
		watcher.Start()
		select {
		case err := <-watcher.ErrorChan:
			assert.Nil(t, err)
		case <-time.After(50 * time.Millisecond):
			return
		}
	})
	t.Run("Start services with error", func(t *testing.T) {
		watcher, err := giles.NewWatcher([]giles.Service{{Name: "foobar", Path: "foobarbaz"}})
		assert.Nil(t, err)
		defer watcher.Close()
		watcher.Start()
		select {
		case err := <-watcher.ErrorChan:
			assert.Error(t, err)
		}
	})
	t.Cleanup(cleanUpTests)
}

func TestWatcher_Watch(t *testing.T) {
	var services []giles.Service
	for _, service := range testServices {
		services = append(services, giles.Service{
			Name: service.name,
			Path: service.path,
		})
	}
	t.Run("Watch files without error", func(t *testing.T) {
		watcher, err := giles.NewWatcher(services)
		assert.Nil(t, err)
		defer watcher.Close()
		err = watcher.AddPaths([]string{filepath.Join(rootPath, ".test")})
		assert.Nil(t, err)
		go watcher.Watch()
		select {
		case err := <-watcher.ErrorChan:
			assert.Nil(t, err)
		case <-time.After(50 * time.Millisecond):
			return
		}
	})
	//TODO Fix failing test
	//t.Run("Watch files with error", func(t *testing.T) {
	//	watcher, err := giles.NewWatcher(services)
	//	assert.Nil(t, err)
	//	defer watcher.Close()
	//	err = watcher.AddPaths([]string{filepath.Join(rootPath, ".test")})
	//	assert.Nil(t, err)
	//	go watcher.Watch()
	//	select {
	//	case err := <-watcher.ErrorChan:
	//		assert.Error(t, err)
	//	}
	//})
	t.Cleanup(cleanUpTests)
}
