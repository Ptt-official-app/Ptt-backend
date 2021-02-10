@echo off
set argC=0
REM calculate number of arguments
for %%x in (%*) do Set /A argC+=1

REM no arguments and number of arguments greater than 1
if %argC% neq 1 goto invalidArgs
if %1==help goto help
if %1==build goto build
if %1==deps goto deps
if %1==format goto format
if %1==lint goto lint
if %1==test-unit goto test-unit
if %1==test-integration goto test-integration
if %1==clean goto clean

:invalidArgs
echo invalid args, please check command
call :help
goto end

:help
echo ---- Project: Ptt-backend ----
echo  Usage: make.bat [COMMAND]
echo.
echo  Management Commands:
echo   build              Build project
echo   deps               Ensures fresh go.mod and go.sum for dependencies
echo   format             Formats Go code
echo   lint               Run golangci-lint check
echo   test-unit          Run all unit tests
echo   test-integration   Run all integration and unit tests
echo   clean              Remove object files, ./bin, .out files
echo.
goto end

REM build: Build project
:build
REM if git tag exists, take it as version, otherwise use commit sha
FOR /F "tokens=*" %%g IN ('"git rev-parse --short HEAD"') do set GITSHA=%%g
set Tag=""
FOR /F %%i IN ('"git rev-list --tags --max-count 1"') do set Tag=%%i
IF NOT %Tag%=="" (
    FOR /F %%j IN ('"git describe --tags %Tag%"') DO set VERSION=%%j
) ELSE (
    set VERSION=git-%GITSHA%
 )

REM get current utc time in format: %Y-%m-%dT%H:%M:%SZ
for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
Set Second=0%Second%
Set Second=%Second:~-2%
Set Minute=0%Minute%
Set Minute=%Minute:~-2%
Set Hour=0%Hour%
Set Hour=%Hour:~-2%
Set Day=0%Day%
Set Day=%Day:~-2%
Set Month=0%Month%
Set Month=%Month:~-2%
set BUILDTIME=%Year%-%Month%-%Day%T%Hour%:%Minute%:%Second%Z

set GOFLAGS=-trimpath
set LDFLAGS="-X main/version.version=%VERSION% -X main/version.commit=%GITSHA% -X main/version.buildTime=%BUILDTIME%"
mkdir bin 2>nul
echo VERSION: %VERSION%
echo GITSHA: %GITSHA%
go build %GOFLAGS% -ldflags %LDFLAGS%
echo executable file .\Ptt-backend.exe
goto end

REM deps: Ensures fresh go.mod and go.sum for dependencies
:deps
go mod tidy
go mod verify
goto end

REM format: Formats Go code
:format
go fmt .\...
goto end

REM lint: Run golangci-lint check
:lint
setlocal
    set "GOBIN=%GOPATH%\bin"
    if not "%GOLANGCI_LINT_VERSION%"=="" (
        set "GOLANGCI_LINT_VERSION=@%GOLANGCI_LINT_VERSION%"
    )
    if not exist "%GOBIN%\golangci-lint.exe" (
        go get github.com/golangci/golangci-lint/cmd/golangci-lint%GOLANGCI_LINT_VERSION%
    )
    %GOBIN%\golangci-lint run ./...
endlocal
goto end

REM test-unit: Run all unit tests
:test-unit
setlocal
    set CGO_ENABLED=1 && go test ./... -coverprofile=coverage.out -cover -race
endlocal
goto end

REM test-integration: Run all integration and unit tests
:test-integration
setlocal
    set CGO_ENABLED=1 && go test ./... -race -tags=integration -covermode=atomic -coverprofile=coverage.tmp
    DEL "*.tmp"
endlocal
goto end

REM clean: Remove object files, ./bin, .out .exe files
:clean
go clean -i -x
echo delete Ptt-backend.exe, clean out files
DEL /Q /F /S "*.out" 2>nul
goto end

:end
