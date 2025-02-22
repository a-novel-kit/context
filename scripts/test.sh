#!/bin/bash

KUBE_FILE="pod.test.yaml"
PG_VOLUME="47870c341f31b19c21f15098c0c9e2fd67fb67602ce5ae04cc4f867ac6a5b7c8"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# First, we set up a temporary directory to receive the coverage (binary) files...
GOCOVERTMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$GOCOVERTMPDIR"' EXIT

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    podman kube down ${KUBE_FILE}
    podman volume rm "${PG_VOLUME}" -f
}
trap int_handler INT

# Setup test containers.
podman play kube ${KUBE_FILE}

export PORT=8080
export DSN="postgres://test:test@localhost:5432/test?sslmode=disable"

# Clear old coverage files.
if [ -d "$GOCOVERTMPDIR" ]; then rm -Rf $GOCOVERTMPDIR; fi
mkdir $GOCOVERTMPDIR

go run ${TEST_TOOL_PKG} --format pkgname -- \
  -cover -covermode=atomic -v -count=1 \
  $(for mod in $(go list -m); do go list ${mod//$(go list .)/.}/...; done) \
  -args -test.gocoverdir=$GOCOVERTMPDIR

# Collect test coverage.
go tool covdata textfmt -i="$GOCOVERTMPDIR" -o=cover.out
go tool cover -html=cover.out -o=cover.html

# Normal execution: containers are shut down.
podman kube down ${KUBE_FILE}
podman volume rm "${PG_VOLUME}" -f
