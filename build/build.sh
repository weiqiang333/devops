#!/bin/bash

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

go build devops.go