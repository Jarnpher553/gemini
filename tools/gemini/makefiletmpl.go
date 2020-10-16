package main

const makefileTmpl = `
RELEASE:=1.0.0
COMMIT:=$(shell git rev-parse --short HEAD)
BUILDTIME:=$(shell date '+%Y-%m-%d %H:%M:%S')
PROJECT:={{name}}
build:
	@GOOS=linux go build -ldflags "-w -s -X '${PROJECT}/version.Commit=${COMMIT}' -X '${PROJECT}/version.BuildTime=${BUILDTIME}' -X '${PROJECT}/version.Release=${RELEASE}'" -o ./{{name}}
clean:
	@rm ./{{name}}
`
