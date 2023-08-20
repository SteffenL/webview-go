@echo off
setlocal

echo Prepare directories...
set script_dir=%~dp0
set src_dir=%script_dir%..
set build_dir=%script_dir%..\build
mkdir "%build_dir%"

echo Webview directory: %src_dir%
echo Build directory: %build_dir%

echo Building Go examples
mkdir build\examples\go
set "CGO_CPPFLAGS=-I%script_dir%\microsoft.web.webview2.%nuget_version%\build\native\include"
set "CGO_LDFLAGS=-L%script_dir%\microsoft.web.webview2.%nuget_version%\build\native\x64"
go build -ldflags="-H windowsgui" -o build\examples\go\basic.exe examples\basic.go || exit /b
go build -ldflags="-H windowsgui" -o build\examples\go\bind.exe examples\bind.go || exit /b

echo Running Go tests
cd /D %src_dir%
set CGO_ENABLED=1
set "PATH=%PATH%;%cd%\build"
go test || exit \b
