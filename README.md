# giles

A watcher of files and runner of multiple Go bianries

I had the idea to create this to be like [nodemon] for multiple Go services running within a monorepo.

## Instructions
```go
package main

import (
  "fmt"
  "github.com/cpustejovsky/giles"
  "os"
  "os/signal"
  "path/filepath"
  "syscall"
)

func main() {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  root := filepath.Join(os.Getenv("YOUR_LOCAL_ROOT_PATH"))
  w, err := giles.NewWatcher(root)
  if err != nil {
    fmt.Println("Error\t", err)
    os.Exit(1)
  }
  defer w.Close()
  //Add services you want to start and restart based on file changes. 
  //the Name property will help you troubleshoot which service is having problems if giles encounters an error
  //the Path property tells giles what go file to build; it currently must point directly to the main.go file
  services := []giles.Service{{
    Name: "service_one",
    Path: "path/to/service_one",
  }, {
    Name: "service_two",
    Path: "path/to/service_two",
  }, {
    Name: "service_three",
    Path: "path/to/service_three",
  }}

  //Tell giles which paths to watch for file changes
  pathOne := filepath.Join(root, "path/to/code/one")
  pathTwo := filepath.Join(root, "path/to/code/two")
  paths := []string{pathOne, pathTwo}
  err = w.AddPaths(paths)
  if err != nil {
    fmt.Println("Error\t", err)
    os.Exit(1)
  }  
  //Start watching for file changes
  go w.Watch(services)
  //Start services
  err = w.Start(services)
  if err != nil {
    err = os.RemoveAll(filepath.Join(root, "tmp/builds"))
    if err != nil {
      fmt.Println("Error\t", err)
      os.Exit(1)
    }
    fmt.Println("Error\t", err)
    os.Exit(1)
  }

  select {
  case sig := <-sigs:
    err = w.Stop()
    if err != nil {
      fmt.Println("Error\t", err)
      os.Exit(1)
    }
    err = os.RemoveAll(filepath.Join(root, "tmp/builds"))
    if err != nil {
      fmt.Println("Error\t", err)
      os.Exit(1)
    }
    fmt.Println("Signal\t", sig)
    os.Exit(0)
  }

}
```

## Next Steps
* Add detailed instructions
* Make paths less brittle. 
  * If the `main.go` file resides inside `~/foo/bar/baz/`, giles should be able to run when it is only given `~/foo`
  * Have it return error if multiple `main.go`s are found in same path with suggestion to break up services
* Add short-circuit for too many file changes (like when a different branch is checked out in git)
  * Potentially having a configurable time to wait before restarting again
* Graceful shutdown within giles
  * Handle tmp file cleanup within giles
* Configurable restarts based on which files changes
  * If the go binary is in the same directory as a file change
* 