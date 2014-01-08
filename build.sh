#!/bin/bash
# check if pwd is in the gopath

echo
echo "Adding current directory to GOPATH: $(pwd)"
# TODO: Check if pwd is in GOPATH and if not, do this command.. hardcoded for now.. but not portable
export GOPATH=$(pwd)
echo "Building binary"
go build -o bin/api_server api_server.go 
echo "Done"
