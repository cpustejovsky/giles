# giles
The magical file watcher for ALL your local Go services.
![giles](./giles.jpeg)
I had the idea to create this to be like [nodemon](https://www.npmjs.com/package/nodemon) for multiple Go services running within a monorepo.
The ultimate goal would be for this to be a CLI like nodemon and to read off a configuration file you point it to.

## Instructions
```go
package main

import (
  "log"
  "github.com/cpustejovsky/giles"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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
  
  w, err := giles.NewWatcher(services)
  if err != nil {
    log.Println("Error\t", err)
    os.Exit(1)
  }
  defer w.Close()

  //Tell giles which paths to watch for file changes
  paths := []string{"path/to/code/one", "path/to/code/two"}
  err = w.AddPaths(paths)
  if err != nil {
    log.Println("Error\t", err)
    os.Exit(1)
  }  
  //Start watching for file changes
  go w.Watch()
  //Start services
  w.Start()

  select {
  case err := <-w.ErrorChan:
    if err != nil {
      log.Println("Error\t", err)
      os.Exit(1)
    }
  case sig := <-sigs:
    log.Println("Signal\t", sig)
    os.Exit(0)
  }

}
```

## Next Steps
* ~~Add detailed instructions~~
* Add error handling to:
  * ~~Start~~
  * Watch
    * Make Watch a testable method
* Make paths less brittle. 
  * If the `main.go` file resides inside `~/foo/bar/baz/`, giles should be able to run when it is only given `~/foo`
  * Have watcher return sentinel error if multiple `main.go`s are found in same path with suggestion to break up services
* Add short-circuit for too many file changes (like when a different branch is checked out in git)
  * Potentially having a configurable time to wait before restarting again
* ~~Graceful shutdown within giles~~
  * ~~Handle tmp file cleanup within giles~~
* Configurable restarts based on which files changes
  * If the go binary is in the same directory as a file change