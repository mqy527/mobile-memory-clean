package main

import (
	"flag"
	"strings"
	"time"

	"github.com/mqy527/mobile-memory-clean/log"
	"github.com/mqy527/mobile-memory-clean/memory/decide"
	"github.com/mqy527/mobile-memory-clean/memory/release"
)

var logger = log.GetLogger("main")

func main() {

	var whitePackages string
	var checkIntervalSecs int64
	var memUsedThreshold float64
	var whitePackagesArr []string

	flag.StringVar(&whitePackages, "whitePackages", "com.android.systemui", "包名白名单（多个包名时，以英文逗号分隔）")
	flag.Int64Var(&checkIntervalSecs, "checkIntervalSecs", 10, "内存监控间隔，单位：秒")
	flag.Float64Var(&memUsedThreshold, "memUsedThreshold", 75, "内存使用率大于该值时，清理内存")
	flag.Parse()

	whitePackages = strings.ReplaceAll(whitePackages, " ", "")
	whitePackages = strings.ReplaceAll(whitePackages, "，", ",")
	whitePackagesArr = strings.Split(whitePackages, ",")

	logger.Infof("checkIntervalSecs: %d s", checkIntervalSecs)
	logger.Infof("memUsedThreshold: %f", memUsedThreshold)
	logger.Infof("whitePackages: %s", whitePackagesArr)

	memReleaseCondition := release.NewMemReleaseConditionByThreshold(memUsedThreshold, memUsedThreshold-5)
	memReleaseDecider := decide.NewWhiteListMemReleaseDecider(whitePackagesArr)
	memReleaser := release.NewMemReleaser(memReleaseCondition, memReleaseDecider)
	for {
		checkAndReleaseMem(memReleaseCondition, memReleaser)
		time.Sleep(time.Second * time.Duration(checkIntervalSecs))
	}
}

func checkAndReleaseMem(memReleaseCondition release.MemReleaseCondition, memReleaser release.MemReleaser) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("error occurs: ", err)
		}
	}()
	if memReleaseCondition.NeedToReleaseNow() {
		logger.Info("begin release memory.")
		err := memReleaser.DoRelease()
		if err != nil {
			logger.Error("DoRelease error: ", err)
		}
	}
}
