@echo off
setlocal
set "GOCACHE=%CD%\.gocache"
set "GOMODCACHE=%CD%\.gomodcache"
if not "%~1"=="" set "HTTP_PORT=%~1"
go run ./cmd/api
