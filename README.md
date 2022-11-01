# giles
The magical file watcher for ALL your local Go services.
![giles](./giles.jpeg)
I had the idea to create this to be like [nodemon](https://www.npmjs.com/package/nodemon) for multiple Go services running within a monorepo.
The ultimate goal would be for this to be a CLI like nodemon and to read off a configuration file you point it to.

## Instructions
First create a yaml file like so:
```yaml
services:
  - name: one
    path: /home/cpustejovsky/go/src/giles/test/one
  - name: two
    path: /home/cpustejovsky/go/src/giles/test/two
  - name: three
    path: /home/cpustejovsky/go/src/giles/test/three
paths:
  - /home/cpustejovsky/go/src/giles/test
  - /home/cpustejovsky/go/src/franz
```

Then set up a main file to run the watcher
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
  //pass in path to configuration yaml file
  w, err := giles.NewWatcher("path/to/your/config.yaml")
  if err != nil {
    log.Println("Error\t", err)
    os.Exit(1)
  }
  defer w.Close()
  
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
* ~~Graceful shutdown within giles~~
  * ~~Handle tmp file cleanup within giles~~
* ~~Read from yaml config file for services and paths to watch~~
* Add finer-grained error handling
  * Specific errors from read method
* Add error handling to:
  * ~~NewWatcher~~
  * ~~Start~~
  * Watch
    * Make Watch a testable method
* Add short-circuit for too many file changes (like when a different branch is checked out in git)
  * Potentially having a configurable time to wait before restarting again

