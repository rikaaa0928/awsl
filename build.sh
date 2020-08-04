#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o awsl.exe
GOOS=linux GOARCH=amd64 go build -o awsl
GOOS=darwin GOARCH=amd64 go build -o awsl_osx
