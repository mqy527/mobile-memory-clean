package command

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mqy527/mobile-memory-clean/log"
)

var logger = log.GetLogger("command")

type PsResult struct {
	USER  string
	PID   string
	PPID  string
	VSIZE string
	RSS   string
	NAME  string
	STIME time.Time
}

func (pr PsResult) String() string {
	return pr.NAME
}

// ListProcessesOrderByRss 按照RSS降序，获取topN个普通用户进程
func ListProcessesOrderByRss(topN int) ([]PsResult, error) {
	cmd := fmt.Sprintf("ps|head -1;ps|grep -v PID|sort -rn -k +5|grep u0_|head -n %d", topN)
	res, err := ExeBash(cmd)
	if err != nil {
		return nil, err
	}
	ret, err := ToPsResult(string(res))
	return ret, err
}

// ListProcessesOrderToyBoxByRss 按照RSS降序，获取topN个普通用户进程
func ListProcessesOrderToyBoxByRss(topN int) ([]PsResult, error) {
	cmd := fmt.Sprintf("toybox ps -A -k -RSS -o USER,PID,PPID,VSZ,RSS,NAME:64,STIME:19|grep u0_|head -n %d", topN)
	res, err := ExeBash(cmd)
	if err != nil {
		return nil, err
	}
	ret, err := ToToyBoxPsResult(string(res))
	return ret, err
}

// ListProcessesByName 列出进程名包含name的所有进程信息。
func ListProcessesByName(name string) ([]PsResult, error) {
	cmd := fmt.Sprintf("toybox ps -A -o USER,PID,PPID,VSZ,RSS,NAME:64,STIME:19|grep %s", name)
	res, err := ExeBash(cmd)
	if err != nil {
		return nil, err
	}
	ret, err := ToToyBoxPsResult(string(res))
	return ret, err
}

// ListProcessByName 列出进程名等于name的进程信息。
func ListProcessByName(name string) (PsResult, error) {
	ret, err := ListProcessesByName(name)
	if err != nil {
		return PsResult{}, err
	}
	for i := 0; i < len(ret); i++ {
		if ret[i].NAME == name {
			return ret[i], nil
		}
	}
	return PsResult{}, fmt.Errorf("process not found by: %s", name)
}

func ToPsResult(res string) ([]PsResult, error) {
	lines := strings.Split(res, "\n")
	if len(lines) <= 0 {
		return nil, nil
	}
	result := make([]PsResult, 0, len(lines))
	for _, element := range lines {
		if len(element) == 0 || strings.Contains(element, "Broken pipe") {
			continue
		}
		logger.Debug(element)
		r := regexp.MustCompile("[^\\s]+")
		eles := r.FindAllString(element, -1)
		result = append(result, PsResult{
			USER:  eles[0],
			PID:   eles[1],
			PPID:  eles[2],
			VSIZE: eles[3],
			RSS:   eles[4],
			NAME:  eles[8],
		})
	}
	return result, nil
}

func ToToyBoxPsResult(res string) ([]PsResult, error) {
	lines := strings.Split(res, "\n")
	if len(lines) <= 0 {
		return nil, nil
	}
	result := make([]PsResult, 0, len(lines))
	for _, element := range lines {
		if len(element) == 0 || strings.Contains(element, "Broken pipe") {
			continue
		}
		logger.Debug(element)
		r := regexp.MustCompile("[^\\s]+")
		eles := r.FindAllString(element, -1)
		stime, _ := time.ParseInLocation("2006-01-02 15:04:05", eles[6]+" "+eles[7], time.Local)
		result = append(result, PsResult{
			USER:  eles[0],
			PID:   eles[1],
			PPID:  eles[2],
			VSIZE: eles[3],
			RSS:   eles[4],
			NAME:  eles[5],
			STIME: stime,
		})
	}
	return result, nil
}
