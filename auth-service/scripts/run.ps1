param(
    [string]$Port
)

$ErrorActionPreference = "Stop"

$env:GOCACHE = Join-Path (Get-Location) ".gocache"
$env:GOMODCACHE = Join-Path (Get-Location) ".gomodcache"
if ($Port) {
    $env:HTTP_PORT = $Port
}

go run ./cmd/api
