@ECHO OFF

CD ..\cmd\create-gp4

REM Windows
SET GOOS=windows
IF NOT EXIST "%OO_PS4_TOOLCHAIN%\bin\windows" MKDIR %OO_PS4_TOOLCHAIN%\bin\windows
go.exe build -o create-gp4.exe
COPY /Y create-gp4.exe %OO_PS4_TOOLCHAIN%\bin\windows\create-gp4.exe
DEL create-gp4.exe

REM Linux
SET GOOS=linux
IF NOT EXIST "%OO_PS4_TOOLCHAIN%\bin\linux" MKDIR %OO_PS4_TOOLCHAIN%\bin\linux
go.exe build -o create-gp4
COPY /Y create-gp4 %OO_PS4_TOOLCHAIN%\bin\linux\create-gp4
DEL create-gp4

REM MacOS
SET GOOS=darwin
IF NOT EXIST "%OO_PS4_TOOLCHAIN%\bin\macos" MKDIR %OO_PS4_TOOLCHAIN%\bin\macos
go.exe build -o create-gp4osx
COPY /Y create-gp4osx %OO_PS4_TOOLCHAIN%\bin\macos\create-gp4
DEL create-gp4osx