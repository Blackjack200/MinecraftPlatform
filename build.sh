#!/bin/bash

function build() {
  echo "[+] Build GOOS=${1} ARCH=${2}"
  CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -o "platform_${1}_${2}" main.go
  mv "platform_${1}_${2}" bin/"platform_${1}_${2}"
}

rm -rdf bin 2>/dev/null >/dev/null
mkdir bin 2>/dev/null >/dev/null

build linux amd64
build linux arm64
build linux arm
build android arm64
build darwin amd64
build windows amd64
