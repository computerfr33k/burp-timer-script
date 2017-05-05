#!/usr/bin/env bash

PLATFORMS=(
    dragonfly-amd64
    freebsd-amd64 freebsd-386
    linux-amd64 linux-386 linux-arm linux-arm64 linux-ppc64 linux-ppc64le linux-mips linux-mipsle
    netbsd-amd64 netbsd-386
    openbsd-amd64 openbsd-386
)

mkdir -p build
echo Building
for plat in "${PLATFORMS[@]}"; do
    echo Building "$plat"

    export GOOS="${plat%-*}"
    export GOARCH="${plat#*-}"
    go build -o "build/timer_script-${plat}" main.go
    echo

done
