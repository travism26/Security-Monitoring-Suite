# External Simulation Design Document (This is the one we are using / building)

### Objective

Build an external simulation tool that generates synthetic metrics resembling threats and sends them directly to the API Gateway for testing the threat detection pipeline.

### Components

1. Simulation Tool
   - Language: Go, Python, or any other language that can send HTTP POST requests.
   - Responsibilities:
     - Generate pre-packaged anomalous data (e.g., high CPU usage, malicious processes).
     - Send data directly to the API Gateway via HTTP POST requests.

Example Payload:

```json
{
  "timestamp": "2024-11-06T12:00:00Z",
  "cpu_usage": 95.0,
  "memory_usage": {
    "used": 16000,
    "total": 32000,
    "percent": 50.0
  },
  "processes": [
    { "name": "malicious.exe", "cpu_percent": 80.0, "memory_percent": 20.0 }
  ]
}
```

Example Code (Go):

```go
   package main

    import (
        "bytes"
        "encoding/json"
        "log"
        "net/http"
        "time"
    )

    type Metrics struct {
        Timestamp   string  `json:"timestamp"`
        CPUUsage    float64 `json:"cpu_usage"`
        MemoryUsage struct {
            Used    int     `json:"used"`
            Total   int     `json:"total"`
            Percent float64 `json:"percent"`
        } `json:"memory_usage"`
        Processes []struct {
            Name         string  `json:"name"`
            CPUPercent   float64 `json:"cpu_percent"`
            MemoryPercent float64 `json:"memory_percent"`
        } `json:"processes"`
    }

    func main() {
        metrics := Metrics{
            Timestamp: time.Now().Format(time.RFC3339),
            CPUUsage:  95.0,
            MemoryUsage: struct {
                Used    int     `json:"used"`
                Total   int     `json:"total"`
                Percent float64 `json:"percent"`
            }{Used: 16000, Total: 32000, Percent: 50.0},
            Processes: []struct {
                Name         string  `json:"name"`
                CPUPercent   float64 `json:"cpu_percent"`
                MemoryPercent float64 `json:"memory_percent"`
            }{
                {Name: "malicious.exe", CPUPercent: 80.0, MemoryPercent: 20.0},
            },
        }

        data, err := json.Marshal(metrics)
        if err != nil {
            log.Fatalf("Error marshalling metrics: %v", err)
        }

        resp, err := http.Post("http://api-gateway-endpoint/metrics", "application/json", bytes.NewBuffer(data))
        if err != nil {
            log.Fatalf("Error sending data: %v", err)
        }
        defer resp.Body.Close()

        log.Printf("Metrics sent with status code: %d", resp.StatusCode)
    }
```

2. API Gateway Interaction
   - The simulation tool sends synthetic payloads to the existing API Gateway endpoint.
   - The API Gateway processes the data and forwards it to Kafka for threat detection analysis.

---

### Design Considerations

- Isolation: Ensure the simulation tool is separate from the production agent.
- Data Accuracy: Simulated data should closely resemble real-world metrics.
- Ease of Use: Provide configuration options for customizing simulated payloads.

---

### Future Enhancements

- Add a GUI or CLI for selecting different simulation scenarios.
- Include invalid/malformed payload testing for error handling validation.

---

## Internal Simulation Design Document

### Objective

Enhance the agent to support a simulation mode that injects synthetic metrics into the data pipeline for testing threat detection capabilities.

### Components

1. Simulation Module
   - Language: Go
   - Responsibilities:
     - Generate synthetic data for CPU, memory, processes, etc.
     - Integrate with the agent’s existing metric collection system.

Example Code:

```go
package collector

import (
    "math/rand"
)

// SimulateCPUUsage generates synthetic CPU usage data
func SimulateCPUUsage() float64 {
    return 95.0 + rand.Float64()*5.0 // Simulates 95%-100% CPU usage
}

// CollectMetrics collects real and simulated metrics
func CollectMetrics(simulate bool) Metrics {
    metrics := collectRealMetrics() // Existing function to pull system metrics

    if simulate {
        metrics.CPUUsage = SimulateCPUUsage()
    }

    return metrics
}

```

2. Configuration
   - Add a simulate_mode flag in the agent’s configuration file (config.yaml):

```yaml
simulate_mode: true
simulate:
cpu_usage: 95.0
process_name: 'malicious.exe'
```

- Load the configuration during the agent’s startup:

```go
package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

type SimulationConfig struct {
    SimulateMode bool    `yaml:"simulate_mode"`
    CPUUsage     float64 `yaml:"cpu_usage"`
    ProcessName  string  `yaml:"process_name"`
}

func LoadConfig(filename string) SimulationConfig {
    var config SimulationConfig
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        log.Fatalf("Error parsing config file: %v", err)
    }
    return config
}

```

3. Data Pipeline
   - Simulated metrics are merged with real metrics during collection.
   - The agent exports the combined data through existing exporters (HTTP/File).

### Design Considerations

- Safety: Ensure that simulation mode is disabled in production environments.
- Logging: Log all synthetic data clearly to avoid confusion with real metrics.
- Configurability: Allow fine-tuned control of simulation scenarios through the configuration file.
  Future Enhancements
  Add support for more complex threat scenarios (e.g., network anomalies, privilege escalation).
  Include a dashboard for toggling simulation mode and viewing simulated data in real time.
  Let me know if you need further adjustments!
