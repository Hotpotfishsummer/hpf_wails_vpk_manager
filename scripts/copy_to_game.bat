@echo off
chcp 65001 >nul
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0copy_to_game.ps1" %*
if %ERRORLEVEL% neq 0 (
    pause
)
