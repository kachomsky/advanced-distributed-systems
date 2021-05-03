@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /L %%x IN (1,1,6) DO (
    SET "NUM=configFile_4%%x.txt
	start cmd.exe /k "go run echo.go !NUM!"
)