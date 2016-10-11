#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -a -ldflags "-s -w" -o bin/nathttpd_darwin_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags "-s -w" -o bin/nathttpd_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -a -ldflags "-s -w" -o bin/nathttpd_windows_amd64.exe