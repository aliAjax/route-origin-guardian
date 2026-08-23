#!/usr/bin/env sh
set -eu
find . -name '*.go' ! -name '*_test.go' -print0 | xargs -0 wc -l
