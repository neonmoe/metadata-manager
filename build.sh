#!/bin/sh
set -e
go get -d -v # Get deps
go build -v # Build program
go test -v # Run tests
go vet # Run vet (additional checks to testing & compilation)
if [ "$1" = "-run" ]; then
    ./metadata-manager
fi
