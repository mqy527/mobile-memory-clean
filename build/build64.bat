set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=1
set CC=D:\adb\arm-linux-gcc\arrch64\bin\aarch64-linux-gnu-gcc
set CXX=D:\adb\arm-linux-gcc\arrch64\bin\aarch64-linux-gnu-g++
go build -o mobile-memory-clean -a -v ..\main.go
pause
:: 安装adb； usb连接mobile； 打开调试模式；
adb shell mkdir -p /data/local/tmp/mobile-memory/logs
adb shell mkdir -p /data/local/tmp/mobile-memory/goroot/lib/time/
adb push mobile-memory-clean /data/local/tmp/mobile-memory/
adb push start.sh /data/local/tmp/mobile-memory/
adb push stop.sh /data/local/tmp/mobile-memory/
adb push zoneinfo.zip /data/local/tmp/mobile-memory/goroot/lib/time/
adb shell
:: cd /data/local/tmp/mobile-memory/
:: source start.sh