#!/bin/bash

KUBE_FILE="pod.test.yaml"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# First, we set up a temporary directory to receive the coverage (binary) files...
GOCOVERTMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$GOCOVERTMPDIR"' EXIT

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    podman kube down ${KUBE_FILE}
}
trap int_handler INT

# Setup test containers.
podman play kube ${KUBE_FILE}

export PORT=8080
export DSN="postgres://test:test@localhost:5432/test?sslmode=disable"

# Clear old coverage files.
if [ -d "$GOCOVERTMPDIR" ]; then rm -Rf $GOCOVERTMPDIR; fi
mkdir $GOCOVERTMPDIR

# Execute tests.
go run ${TEST_TOOL_PKG} --format pkgname -- \
  -cover -covermode=atomic -v -count=1 \
  $(go list -m | grep -v /mocks) \
  -args -test.gocoverdir=$GOCOVERTMPDIR

# Collect test coverage.
go tool covdata textfmt -i="$GOCOVERTMPDIR" -o=cover.out
go tool cover -html=cover.out -o=cover.html

# Normal execution: containers are shut down.
podman kube down ${KUBE_FILE}
