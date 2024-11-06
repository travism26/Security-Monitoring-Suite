package main

import (
    "github.com/travism26/windows-monitoring-agent/pkg/collector"
    "github.com/travism26/windows-monitoring-agent/pkg/sender"
    "log"
)

func main() {
    log.Println("Starting Windows Monitoring Agent...")

    // Example: Collect CPU usage and send to dashboard
    cpuUsage := collector.CollectCPUUsage()
    sender.SendData(cpuUsage)

    log.Println("Monitoring Agent finished collecting data.")
}
