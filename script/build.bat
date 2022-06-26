@echo off

echo Prepare directories...
set script_dir=%~dp0
set src_dir=%script_dir%..
set build_dir=%script_dir%..\build
mkdir "%build_dir%"

echo Webview directory: %src_dir%
echo Build directory: %build_dir%

:: If you update the nuget package, change its version here
set nuget_version=1.0.1150.38
echo Using Nuget Package microsoft.web.webview2.%nuget_version%
if not exist "%script_dir%\microsoft.web.webview2.%nuget_version%" (
	curl -sSLO https://dist.nuget.org/win-x86-commandline/latest/nuget.exe
	nuget.exe install Microsoft.Web.Webview2 -Version %nuget_version% -OutputDirectory %script_dir%
	echo Nuget package installed
)

if not exist "%build_dir%\WebView2Loader.dll" (
	copy "%script_dir%\microsoft.web.webview2.%nuget_version%\build\native\x64\WebView2Loader.dll" "%build_dir%"
)



echo Building Go examples
mkdir build\examples\go
set "CGO_CPPFLAGS=-I%script_dir%\microsoft.web.webview2.%nuget_version%\build\native\include"
set "CGO_LDFLAGS=-L%script_dir%\microsoft.web.webview2.%nuget_version%\build\native\x64"
go build -ldflags="-H windowsgui" -o build\examples\go\basic.exe examples\basic.go || exit /b
go build -ldflags="-H windowsgui" -o build\examples\go\bind.exe examples\bind.go || exit /b

echo Running Go tests
cd /D %src_dir%
set CGO_ENABLED=1
go test || exit \b
