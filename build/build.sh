#!/bin/bash

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

go build -o bin/devops devops.go
go build -o bin/devops-cron examples/devops-cron/devops-cron.go

zip devops.zip -r bin/ web/
