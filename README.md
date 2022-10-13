# giles

A watcher of files and runner of multiple Go bianries

I had the idea to create this to be like [nodemon] for multiple Go services running within a monorepo.

## Instructions
Coming Soon!

## Next Steps
* Add detailed instructions
* Make paths less brittle. 
  * If the `main.go` file resides inside `~/foo/bar/baz/`, giles should be able to run when it is only given `~/foo`
  * Have it return error if multiple `main.go`s are found in same path with suggestion to break up services
* Add short-circuit for too many file changes (like when a different branch is checked out in git)
  * Potentially having a configurable time to wait before restarting again
* Configurable restarts based on which files changes
  * If the go binary is in the same directory as a file change
* 