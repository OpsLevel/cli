#! /bin/sh

# Stub script to call `go run` but with the right working directory so that go's
# module system is happy, when called by systems like Claude Desktop that don't
# permit setting the working directory. Primarily useful for development.

cd "$(dirname "$0")"
go run main.go "$@"
