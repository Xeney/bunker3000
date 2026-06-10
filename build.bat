@echo off
set PATH=C:\msys64\ucrt64\bin;%PATH%
set CGO_ENABLED=1
go build -o Bunker3000.exe .
if %errorlevel% equ 0 (
    echo.
    echo OK: Bunker3000.exe ready
) else (
    echo.
    echo FAILED
)
pause
