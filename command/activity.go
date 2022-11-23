package command

import (
	"strings"
)

// GetTopActivityPackage 获取用户当前正在操作的应用包名
func GetTopActivityPackage() string {
	cmd := "dumpsys activity top|grep ACTIVITY"
	// ACTIVITY com.clcw.exejia/.activity.WelcomeActivity 5e913a7 pid=18296
	res, err := ExeBash(cmd)
	if err != nil {
		return ""
	}
	return splitTopActivityPackage(string(res))
}

func splitTopActivityPackage(activityInfo string) string {
	activityInfo = strings.TrimSpace(activityInfo)
	resArr := strings.Split(activityInfo, " ")
	if len(resArr) == 4 {
		pkgs := strings.Split(resArr[1], "/")
		return pkgs[0]
	}
	return ""
}


func GetRecentActivityPackages() []string {
	cmd := "dumpsys activity recents|grep realActivity"
	// 		realActivity=com.coolapk.market/.view.main.MainActivity
	// 		realActivity=com.xiaomi.smarthome/.SmartHomeMainActivity
	res, err := ExeBash(cmd)
	if err != nil {
		return nil
	}
	
	lines := strings.Split(string(res), "\n")
	prefix := "realActivity="
	startIndex := len(prefix)
	result := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) <= startIndex || !strings.HasPrefix(line, prefix) {
			continue
		}
		result = append(result, strings.Split(line[startIndex:], "/")[0])
	}
	return result
}