package collector

import (
    "github.com/shirou/gopsutil/cpu"
    "log"
)

// CollectCPUUsage collects CPU usage information.
func CollectCPUUsage() float64 {
    percent, err := cpu.Percent(0, false)
    if err != nil {
        log.Printf("Error collecting CPU usage: %v", err)
        return 0
    }
    return percent[0]
}
