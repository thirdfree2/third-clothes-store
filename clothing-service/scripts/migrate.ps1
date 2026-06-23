param(
    [string]$Command = "up"
)

$ErrorActionPreference = "Stop"

$env:GOCACHE = Join-Path (Get-Location) ".gocache"
$env:GOMODCACHE = Join-Path (Get-Location) ".gomodcache"

go run ./cmd/migrate $Command
