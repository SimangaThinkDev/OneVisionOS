package monitor

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/host"
	"time"
)

type SystemMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Uptime      uint64  `json:"uptime"`
}

func GetSystemMetrics() SystemMetrics {
	var metrics SystemMetrics

	// CPU Usage
	percentages, err := cpu.Percent(time.Second, false)
	if err == nil && len(percentages) > 0 {
		metrics.CPUUsage = percentages[0]
	}

	// Memory Usage
	vMem, err := mem.VirtualMemory()
	if err == nil {
		metrics.MemoryUsage = vMem.UsedPercent
	}

	// Disk Usage
	dUsage, err := disk.Usage("/")
	if err == nil {
		metrics.DiskUsage = dUsage.UsedPercent
	}

	// Uptime
	uTime, err := host.Uptime()
	if err == nil {
		metrics.Uptime = uTime
	}

	return metrics
}

func CalculateSecurityScore(nidsEngineSignatures int, alertsCount int) int {
	// Simple formula: Start at 100, deduct for alerts, reward for active signatures
	score := 100
	
	// Deduct for alerts (Phase 14 will bring real alert counts)
	// For now, let's assume if we have many signatures, it's safer
	if nidsEngineSignatures < 5 {
		score -= 10
	}
	
	// Ensure score is between 0 and 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	return score
}
