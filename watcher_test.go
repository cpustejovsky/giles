package giles_test

import (
	"fmt"
	"github.com/cpustejovsky/giles"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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

func setUpTests() error {
	for _, service := range testServices {
		err := exec.Command(filepath.Join(rootPath, "test_setup.sh"), service.name, service.port).Run()
		if err != nil {
			return err
		}
		service.path = filepath.Join(rootPath, fmt.Sprintf("/test/%v/main.go", service.name))
	}
	return nil
}

func cleanUpTestDirectory() {
	err := os.RemoveAll(filepath.Join(rootPath, "test"))
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
	err = setUpTests()
	if err != nil {
		log.Fatal(err)
	}
}

func TestWatcher_AddPaths(t *testing.T) {
	err := setUpTests()
	assert.Nil(t, err)
	watcher, err := giles.NewWatcher(rootPath)
	assert.Nil(t, err)
	defer watcher.Close()
	t.Run("Add existing directories to watch", func(t *testing.T) {
		err := watcher.AddPaths([]string{filepath.Join(rootPath, "test")})
		assert.Nil(t, err)
	})
	t.Run("Add unknown directories to watch", func(t *testing.T) {
		err := watcher.AddPaths([]string{filepath.Join(rootPath, "foobarbaz")})
		assert.Error(t, err)
	})
	t.Cleanup(cleanUpTestDirectory)
}

func TestWatcher_Start(t *testing.T) {
	err := setUpTests()
	assert.Nil(t, err)
	watcher, err := giles.NewWatcher(rootPath)
	assert.Nil(t, err)
	defer watcher.Close()
	t.Run("Run services without error", func(t *testing.T) {
		var services []giles.Service
		for _, service := range testServices {
			services = append(services, giles.Service{
				Name: service.name,
				Path: service.path,
			})
		}
		err := watcher.Start(services)
		assert.Nil(t, err)
	})
	t.Run("Run services with error", func(t *testing.T) {
		err := watcher.Start([]giles.Service{{Name: "foobar", Path: "foobarbaz"}})
		assert.Nil(t, err)
	})
	t.Cleanup(cleanUpTestDirectory)
}
