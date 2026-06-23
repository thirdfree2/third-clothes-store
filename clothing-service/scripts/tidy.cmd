@echo off
setlocal
set "GOCACHE=%CD%\.gocache"
set "GOMODCACHE=%CD%\.gomodcache"
go mod tidy
