package giles_test

import (
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
