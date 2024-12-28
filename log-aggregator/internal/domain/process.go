package domain

type Process struct {
	ID          string  `json:"id"`
	LogID       string  `json:"log_id"`
	Name        string  `json:"name"`
	PID         int     `json:"pid"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsage int64   `json:"memory_usage"`
	Status      string  `json:"status"`
}
