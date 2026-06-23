@echo off
setlocal

set SWAG=%USERPROFILE%\go\bin\swag.exe

echo Generating Swagger docs...
"%SWAG%" init -g cmd/api/main.go -o docs --parseDependency --parseInternal

echo Done.
endlocal