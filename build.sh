#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o awsl.exe
GOOS=linux GOARCH=amd64 go build -o awsl
GOOS=darwin GOARCH=amd64 go build -o awsl_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o awsl_darwin_arm64
