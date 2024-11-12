package monitor

import (
	"fmt"
	"runtime"

	"github.com/travism26/system-monitoring-agent/internal/os_darwin"
	"github.com/travism26/system-monitoring-agent/internal/os_windows"

	"github.com/travism26/system-monitoring-agent/internal/core"
)

func NewSystemMonitor() (core.Monitor, error) {
	switch runtime.GOOS {
	case "darwin":
		return os_darwin.NewDarwinMonitor(), nil
	case "windows":
		return os_windows.NewWindowsMonitor(), nil
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
