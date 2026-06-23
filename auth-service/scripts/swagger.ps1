param()

$ErrorActionPreference = "Stop"

$env:GOCACHE = Join-Path (Get-Location) ".gocache"
$env:GOMODCACHE = Join-Path (Get-Location) ".gomodcache"

go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go --parseDependency --parseInternal
