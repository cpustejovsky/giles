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
    path: /home/user/go/src/giles/test/one
  - name: two
    path: /home/user/go/src/giles/test/two
  - name: three
    path: /home/user/go/src/giles/test/three
paths:
  - /home/user/go/src/giles/test
  - /home/user/go/src/franz
```

Then download the package and, inside the gile repository, run:
```shell
go run ./cmd/main.go /path/to/config.yaml
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
* Make sure logging begins when watcher starts, not just on restart
