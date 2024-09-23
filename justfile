set positional-arguments
set quiet
set dotenv-load

default:
  just --list

[group('build')]
build: build-client build-server

[group('build')]
build-client:
    go build -o ./out/client ./bin/client

[group('build')]
build-server:
    go build -o ./out/server ./bin/server

[no-exit-message]
[group('runtime')]
cmd *args:
    go run ./bin/client "$@"

[group('runtime')]
run:
    go run ./bin/server

[group('dev')]
fmt:
    go fmt ./bin/server ./bin/client ./pkgs/types
