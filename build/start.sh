#!/system/bin/sh
export PATH="${PWD}":"$PATH"
export GOROOT="${PWD}"/goroot
export TZ=Asia/Shanghai
export logLevel=info
export logPath="${PWD}"/logs
chmod 755 mobile-memory-clean
#-checkIntervalSecs 内存使用率检测时间间隔，单位：秒。
#-memUsedThreshold 内存使用率上限，内存使用率达到该值时，开始清理内存。
#-whitePackages 白名单，包名，多个包名以逗号分隔。不清理白名单内的进程。
nohup mobile-memory-clean -checkIntervalSecs=3 -memUsedThreshold=75 -whitePackages="com.meizu,com.flyme,com.oma.drm.server,com.huawei,com.android,android,com.google.android,app_process,sh,ping,com.example.cpuactive" &