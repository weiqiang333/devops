#!/bin/bash

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

go build devops.go
go build examples/devops-cron/devops-cron.go
