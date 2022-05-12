#!/bin/bash

mkdir -p dist
go build --ldflags="-s -w" -o ./dist/ sshd.go
go build --ldflags="-s -w" -o ./dist/ telnetd.go
go build --ldflags="-s -w" -o ./dist/ fshell.go
