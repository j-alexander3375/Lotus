@echo off
REM Automated Release Wrapper for Windows
REM Usage: release.bat [minor|patch]

set BUMP_TYPE=%1
if "%BUMP_TYPE%"=="" set BUMP_TYPE=patch

echo.
echo ========================================
echo   Lotus Automated Release System
echo ========================================
echo.
echo Bump type: %BUMP_TYPE%
echo.

wsl bash -c "cd /mnt/c/Users/joshu/develLotus && ./scripts/auto_release.sh %BUMP_TYPE%"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo   Release Completed Successfully!
    echo ========================================
    echo.
) else (
    echo.
    echo ========================================
    echo   Release Failed!
    echo ========================================
    echo.
    echo Check logs in test_results.log and fix_log.txt
)

pause
