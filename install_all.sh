#!/bin/bash
GOOS=windows GOARCH=amd64 go install
GOOS=linux GOARCH=amd64 go install