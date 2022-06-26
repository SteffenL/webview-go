#!/bin/sh

set -e

mkdir -p build/examples/go || true

if command -v go >/dev/null 2>&1 ; then
	echo "Running Go tests"
	CGO_ENABLED=1 go test
else
	echo "SKIP: Go tests"
fi
