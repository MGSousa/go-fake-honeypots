#!/bin/bash

mkdir -p dist
go build --ldflags="-s -w" -buildmode=pie -o ./dist/ sshd.go
go build --ldflags="-s -w" -buildmode=pie -o ./dist/ telnetd.go
go build --ldflags="-s -w" -buildmode=pie -o ./dist/ fshell.go
