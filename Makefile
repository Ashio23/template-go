 .PHONY: build gen run exec cpv test

CURRENT_DIR = $(shell pwd)

build: gen
	go build -o ./tmp/main .

gen:
	PATH="$(PATH):$(CURRENT_DIR)/bin" go generate ./...

run:
	bin/CompileDaemon -build="make build" -command="make exec" -color=true -exclude="*_enum.go" -graceful-kill=true

exec:
	./tmp/main

cov:
	go test -cover ./...

test:
	go test ./...

go-mod:
	go mod download
	go mod tidy
	go mod verify

install-go-enum:
	GOBIN=$(CURRENT_DIR)/bin go install github.com/abice/go-enum@latest

install-compile-daemon:
	GOBIN=$(CURRENT_DIR)/bin go install github.com/githubnemo/CompileDaemon@latest

first-install: go-mod install-go-enum install-compile-daemon
 