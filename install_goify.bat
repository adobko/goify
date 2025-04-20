@echo off
echo Compiling Goify...
go build -o goify.exe ./app

if %ERRORLEVEL% neq 0 (
    echo Compilation failed. Please make sure Go is installed and set up correctly.
    pause
    exit /b
)

echo Enter the installation directory (default is C:\Goify):
set /p INSTALL_DIR="Install Path: "
if "%INSTALL_DIR%"=="" set "INSTALL_DIR=C:\Goify"

echo Creating installation folder...
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

echo Moving goify.exe to %INSTALL_DIR%...
move goify.exe "%INSTALL_DIR%"

echo Adding %INSTALL_DIR% to PATH...
setx PATH "%PATH%;%INSTALL_DIR%"

if %ERRORLEVEL% neq 0 (
    echo Failed to update PATH. You may need to do it manually.
    pause
    exit /b
)

echo Installation complete. Goify is ready to use from the terminal!
pause