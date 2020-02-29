#!/bin/bash

APP_NAME=timesheet

go get -u -v github.com/go-bindata/go-bindata/...@v3.1.2

# build assets bindata.go file
$GOPATH/bin/go-bindata -fs ./assets/... ./templates

# build executable for linux, osx and windows
GOOS=linux GOARCH=amd64 go build -o $APP_NAME-linux64
zip -9 $APP_NAME-linux64.zip $APP_NAME-linux64
GOOS=darwin GOARCH=amd64 go build -o $APP_NAME-osx64
zip -9 $APP_NAME-osx64.zip $APP_NAME-osx64
GOOS=windows GOARCH=amd64 go build -o $APP_NAME-win64.exe
zip -9 $APP_NAME-win64.zip $APP_NAME-win64.exe
