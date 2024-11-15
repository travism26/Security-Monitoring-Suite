package utils

import (
	"fmt"
)

const (
	_          = iota // ignore first value by assigning to blank identifier
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
)

func FormatBytes(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case value >= TB:
		unit = "TB"
		value = value / TB
	case value >= GB:
		unit = "GB"
		value = value / GB
	case value >= MB:
		unit = "MB"
		value = value / MB
	case value >= KB:
		unit = "KB"
		value = value / KB
	default:
		unit = "B"
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

func FormatNetworkSpeed(bytesPerSecond float64) string {
	return fmt.Sprintf("%s/s", FormatBytes(uint64(bytesPerSecond)))
}
