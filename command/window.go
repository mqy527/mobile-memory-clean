package command

import (
	"strings"
)

func GetTopWindowPackage() string {
	cmd := "dumpsys window | grep mCurrentFocus"
	// mCurrentFocus=Window{878e21e u0 com.clcw.exejia/com.clcw.exejia.activity.MainActivity}
	res, err := ExeBash(cmd)
	if err != nil {
		return ""
	}
	return splitTopWindowPackage(string(res))
}

func splitTopWindowPackage(windowInfo string) string {
	windowInfo = strings.TrimSpace(windowInfo)
	resArr := strings.Split(windowInfo, " ")
	if len(resArr) == 3 {
		pkgs := strings.Split(resArr[2], "/")
		return pkgs[0]
	}
	return ""
}
