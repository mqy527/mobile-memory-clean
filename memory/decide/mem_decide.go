package decide

import (
	"sort"
	"strings"

	"github.com/mqy527/mobile-memory-clean/command"
	"github.com/mqy527/mobile-memory-clean/log"
)

var logger = log.GetLogger("mem-decide")

// MemReleaseDecider 决定结束哪些进程
type MemReleaseDecider interface {
	// DecidedToRelease获取需要被结束的进程，有序地结束进程。
	DecidedToRelease() []command.PsResult
}

// NewWhiteListMemReleaseDecider 基于应用白名单，来决定结束哪些进程
func NewWhiteListMemReleaseDecider(whitePackages []string) MemReleaseDecider {
	mrd := whiteListMemReleaseDecider{
		whitePackages: whitePackages,
	}
	return mrd
}

type whiteListMemReleaseDecider struct {
	//whitePackages 白名单内的，不会被结束
	whitePackages []string
}

func (mrd whiteListMemReleaseDecider) DecidedToRelease() []command.PsResult {
	ps, err := command.ListProcessesOrderToyBoxByRss(50)
	if err != nil {
		logger.Errorf("acquire process info error: %v", err)
		return nil
	}
	ret := make([]command.PsResult, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		if mrd.canRelease(ps[i]) {
			ret = append(ret, ps[i])
		}
	}
	sortByStartTimeAsc(ret)
	ret = sortByRecentUse(ret)
	return ret
}

func (mrd whiteListMemReleaseDecider) canRelease(psResult command.PsResult) bool {
	if strings.Contains(psResult.NAME, ":") || strings.Contains(psResult.NAME, "[") || strings.Contains(psResult.NAME, "/") {
		return false
	}
	if mrd.inWhiteList(psResult.NAME) {
		logger.Debugf("[%s] is in white list, will not be killed.", psResult.NAME)
		return false
	}
	return true
}

func (mrd whiteListMemReleaseDecider) inWhiteList(pkgName string) bool {
	for i := 0; i < len(mrd.whitePackages); i++ {
		if strings.HasPrefix(pkgName, mrd.whitePackages[i]) {
			return true
		}
	}
	return false
}

// sortByStartTimeAsc 对ps，按照启动时间排序
func sortByStartTimeAsc(ps []command.PsResult) {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].STIME.Before(ps[j].STIME)
	})
}

// sortByRecent 对ps，按照最近使用情况排序，最近使用的排后面
func sortByRecentUse(ps []command.PsResult) []command.PsResult {
	recentPackages := command.GetRecentActivityPackages()
	return sortByRecentUsePkgs(ps, recentPackages)
}

func sortByRecentUsePkgs(ps []command.PsResult, recentPackages []string) []command.PsResult {
	EMPTY := command.PsResult{NAME: "MY_EMPYT_001"}
	result := make([]command.PsResult, len(ps))
	resultIndex := len(ps)

	for i := 0; i < len(recentPackages); i++ {
		for j := 0; j < len(ps); j++ {
			if strings.HasPrefix(ps[j].NAME, recentPackages[i]) {
				resultIndex--
				result[resultIndex] = ps[j]
				ps[j] = EMPTY
			}
		}
	}
	i := 0
	for j := 0; j < len(ps); j++ {
		if ps[j] != EMPTY {
			result[i] = ps[j]
			i++
		}
	}

	return result
}
