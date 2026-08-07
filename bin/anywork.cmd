@echo off
set "ANYWORK_REPO=%~dp0.."
set "PYTHONPATH=%ANYWORK_REPO%\src;%PYTHONPATH%"
py -3 -m anywork %*
