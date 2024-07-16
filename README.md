>Helping binaries

- go-enum: let us build enums that helps us validate data with a set of predetermined values. This library/binary creates go files with "go generate" commands, so the binary itself is not needed for the actual Go code, it can be gitignored.
- CompileDaemon: helps rebuilding the code each time we apply a change on any .go files. It is configured so that the "go generate" commands are executed each time it rebuilds. This also kills any http request open and kills the whole process. This is useful only for development so it can be gitignored.

>Commands

Run following commands on the root folder for the project.

>First install:
```
make first-install
```
Installs go dependencies, go-enum (for enums) and CompileDaemon (live reload). The 2 resulting binaries are not needed for the production code, as go-enum creates go files with "go generate" commands and CompileDaemon is just used for development.

>Execute
```
make run
```
This runs `go generate ./...`, builds the project's binary and then executes it.
Any changes on .go files triggers another generate/build/run cycle.

## Approval process Pull Request

Each Pull Request that is created must comply with the [template](pull_request_template.md).