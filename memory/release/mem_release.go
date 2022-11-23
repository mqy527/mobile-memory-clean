package release

import (
	"fmt"
	"strings"
	"time"

	"github.com/mqy527/mobile-memory-clean/command"
	"github.com/mqy527/mobile-memory-clean/log"
	"github.com/mqy527/mobile-memory-clean/memory"
	"github.com/mqy527/mobile-memory-clean/memory/decide"
)

var logger = log.GetLogger("mem-release")

// MemReleaser 释放内存
type MemReleaser interface {
	// DoRelease 释放内存。
	DoRelease() error
}

// MemReleaseCondition 确定释放内存的条件
type MemReleaseCondition interface {
	// NeedToReleaseNow 现在是否需要释放内存
	NeedToReleaseNow() bool

	// StopToReleaseNow 现在是否停止释放内存
	StopToReleaseNow() bool
}

// NewMemReleaseConditionByThreshold 根据内存使用阈值来确定是否释放内存。
// 当前内存使用率>=memUsedThresholdMax，则可以释放内存。
// 当前内存使用率<memUsedThresholdMin，则停止释放内存。
func NewMemReleaseConditionByThreshold(memUsedThresholdMax float64, memUsedThresholdMin float64) MemReleaseCondition {
	rc := memReleaseCondition{
		memUsedThresholdMax: memUsedThresholdMax,
		memUsedThresholdMin: memUsedThresholdMin,
	}
	return rc
}

type memReleaseCondition struct {
	memUsedThresholdMax float64
	memUsedThresholdMin float64
}

func (rc memReleaseCondition) NeedToReleaseNow() bool {
	memoryCurrentUsedPercent := memory.GetUsedPercent()
	logger.Infof("current used memory percent:%.2f", memoryCurrentUsedPercent)
	return memoryCurrentUsedPercent >= rc.memUsedThresholdMax
}

func (rc memReleaseCondition) StopToReleaseNow() bool {
	memoryCurrentUsedPercent := memory.GetUsedPercent()
	logger.Infof("current used memory percent:%.2f", memoryCurrentUsedPercent)
	return memoryCurrentUsedPercent < rc.memUsedThresholdMin
}

func NewMemReleaser(mrc MemReleaseCondition, mrd decide.MemReleaseDecider) MemReleaser {
	mr := memReleaser{
		MemReleaseCondition:      mrc,
		MemReleaseDecider:        mrd,
		pauseSecondsWhenConflict: 3,
	}
	return mr
}

type memReleaser struct {
	MemReleaseCondition
	decide.MemReleaseDecider
	pauseSecondsWhenConflict int
}

func (mr memReleaser) DoRelease() error {
	// 0、避免打开app与清理过程冲突
	// latestFocusPkgName是用户当前正在操作的应用包名
	latestFocusPkgName := command.GetTopActivityPackage()
	latestPs, err := command.ListProcessByName(latestFocusPkgName)
	if err == nil && latestPs.STIME.Add(3*time.Second).After(time.Now()) {
		// 如果最近的app刚启动，则暂停一会儿
		logger.Infof("current[%s] process is starting, pause for a moment.", latestFocusPkgName)
		time.Sleep(time.Duration(mr.pauseSecondsWhenConflict) * time.Second)
	}

	// 1、获取需要被清理的进程
	ps := mr.DecidedToRelease()
	logger.Debug("kill order: ", ps)

	// 2、开始清理进程
	for index, psResult := range ps {
		// currentFocusPkgName是用户当前正在操作的应用包名
		currentFocusPkgName := command.GetTopActivityPackage()
		if len(currentFocusPkgName) > 0 {
			if strings.Contains(psResult.NAME, currentFocusPkgName) {
				// 不能清理用户当前正在操作的应用
				logger.Infof("currently focused process[%s], will not be killed.", psResult.NAME)
				latestFocusPkgName = currentFocusPkgName
				continue
			}
			if currentFocusPkgName != latestFocusPkgName {
				//清理过程中，如果打开新的app，暂停一会儿
				logger.Infof("focused process changed([%s]-->[%s]), pause for a moment.", latestFocusPkgName, currentFocusPkgName)
				time.Sleep(time.Duration(mr.pauseSecondsWhenConflict) * time.Second)
			}
			latestFocusPkgName = currentFocusPkgName
		}
		err := mr.kill(psResult)
		if err != nil {
			logger.Errorf("kill [%s] error: %v", psResult.NAME, err)
		} else {
			logger.Infof("killed: %s", psResult.NAME)
		}
		if index < len(ps)-1 && mr.StopToReleaseNow() {
			logger.Infof("end kill early. [%d/%d]", index+1, len(ps))
			break
		}
	}

	return nil
}

func (mr memReleaser) kill(psResult command.PsResult) error {
	killCmd := "am force-stop %s"
	exeKillCmd := fmt.Sprintf(killCmd, psResult.NAME)
	_, err := command.ExeBash(exeKillCmd)
	return err
}
