@echo off
setlocal
set "GOCACHE=%CD%\.gocache"
set "GOMODCACHE=%CD%\.gomodcache"
if "%~1"=="" (
  go run ./cmd/migrate up
) else (
  go run ./cmd/migrate %*
)
