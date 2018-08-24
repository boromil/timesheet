#!/bin/bash

APP_NAME=timesheet

# make sure that govendor is installed
go get -u github.com/kardianos/govendor
# sync dependencies
govendor sync

# build assets bindata.go file
go-bindata ./assets/... ./templates/index.html

# build executable for linux, osx and windows
GOOS=linux GOARCH=amd64 go build -o $APP_NAME-linux64
zip -9 $APP_NAME-linux64.zip $APP_NAME-linux64
GOOS=darwin GOARCH=amd64 go build -o $APP_NAME-osx64
zip -9 $APP_NAME-osx64.zip $APP_NAME-osx64
GOOS=windows GOARCH=amd64 go build -o $APP_NAME-win64.exe
zip -9 $APP_NAME-win64.zip $APP_NAME-win64.exe
