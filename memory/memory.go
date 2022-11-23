package memory

import "github.com/shirou/gopsutil/mem"

func GetUsedPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetAvailableBytes() uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Available
}