@echo off
setlocal

REM Create dist directory
if not exist dist mkdir dist

REM Build for Windows (amd64)
echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o dist\ozul-windows-amd64.exe .

REM Build for Linux (amd64)
echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o dist\ozul-linux-amd64 .

REM Build for macOS (amd64)
echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o dist\ozul-darwin-amd64 .

echo Build complete! Binaries are in the dist\ directory.
endlocal 