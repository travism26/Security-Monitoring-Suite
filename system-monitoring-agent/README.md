# System Monitoring Agent

A cross-platform system monitoring agent built in Go that collects and logs detailed system metrics. Currently supports macOS and Windows.

## Features

- Real-time system metrics collection:
  - CPU usage
  - Memory usage (used, total, percentage)
  - Disk statistics (total, used, free, usage percentage)
  - Network statistics (bytes sent/received, transfer rates)
- JSON-formatted logging for easy integration with log analyzers (Splunk, ELK Stack)
- Configurable monitoring intervals
- Cross-platform support (macOS, Windows)

# Design Ideas

## Architecture

This is a high-level architecture diagram of the system monitoring agent.

Monitoring Agent is a lightweight, cross-platform application that collects system metrics and logs them to a file. I will be deploying all my other applications on Kubernetes, so I want to keep this lightweight and easy to deploy.

```plaintext
Architecture:

IOT Device/Endpoint                  Kubernetes Cluster
+---------------+                    +----------------------------------------+
|               |                    |  +-----------+         +-----------+   |
| Monitoring    |  HTTPS POST/PUT    |  |           |  Kafka  |           |   |
| Agent         | -----------------> |  |   API     | ------> |  Kafka    |   |
|               |                    |  |  Gateway  |         | Cluster   |   |
+---------------+                    |  |           |         |           |   |
                                     |  +-----------+         +-----------+   |
                                     |        ↑                               |
                                     |        | Horizontal Pod Autoscaling    |
                                     +----------------------------------------+
```

## Prerequisites

- Go 1.22 or higher
- Access to system metrics (may require elevated privileges)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/system-monitoring-agent.git
cd system-monitoring-agent
```

2. Install dependencies:

```bash
go mod download
```

## Configuration

The agent can be configured through `configs/config.yaml`:

```yaml
LogFilePath: "./agent.log"
Interval: 10 # Polling interval in seconds
Monitors:
  CPU: true
  Memory: true
  Disk: true
  Network: true
```

## Building

```bash
go build -o monitoring-agent ./cmd/agent/main.go
```

## Running

```bash
./monitoring-agent
```

## Example Output

```json
{
  "timestamp": 1731791114,
  "timestamp_utc": "2024-11-16T21:05:14Z",
  "host_info": {
    "os": "darwin",
    "arch": "arm64"
  },
  "cpu_usage": 9.37,
  "memory_usage": 11153670144,
  "memory_usage_percent": 57.71,
  "total_memory": 19327352832,
  "disk": {
    "free": 286532517888,
    "total": 494384795648,
    "usage_percent": 42.04,
    "used": 207852277760
  },
  "network": {
    "bytes_received": 59426561897,
    "bytes_received_per_second": 352.5,
    "bytes_sent": 18850020613,
    "bytes_sent_per_second": 2274.89
  },
  "threat_indicators": [
    {
      "type": "high_cpu_usage",
      "description": "CPU usage exceeds threshold",
      "severity": "low",
      "score": 11.25,
      "timestamp": "2024-11-16T17:05:14.579362-04:00",
      "metadata": {
        "tags": ["performance", "resource_usage"]
      }
    }
  ],
  "processes": {
    "process_list": [
      {
        "cpu_percent": 0.5087533993160189,
        "memory_usage": 26148864,
        "name": "launchd",
        "pid": 1,
        "status": "S"
      },
      {
        "cpu_percent": 0.2202767781072918,
        "memory_usage": 67502080,
        "name": "logd",
        "pid": 510,
        "status": "S"
      }
      // ... Truncated for brevity
    ],
    "total_count": 589,
    "total_cpu_percent": 79.00824584831197,
    "total_memory_usage": 16922177536
  }
}
```

## Viewing Logs

```bash
cat agent.log
```

## Development Status

This project is currently under active development. Upcoming features:

- [x] Process monitoring (Completed 2024-11-22)
- [ ] Push payloads to API Gateway (http / rest) which will be sent to Kafka
- [ ] System temperature monitoring
- [ ] Alert configurations (Maybe not this should be collection only alerts will be in SIEM / XDR side)
- [ ] Web interface for metrics visualization (MAYBE NOT KEEP THIS APP SMALL)
- [ ] Docker containerization
- [ ] Support for Linux systems

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Known Issues

- IOMasterPort deprecation warning on macOS 12+ (from gopsutil dependency)
- Limited Windows network statistics

## License

Apache License 2.0 (see LICENSE file)

## Acknowledgments

- [gopsutil](https://github.com/shirou/gopsutil) for system metrics collection
- [viper](https://github.com/spf13/viper) for configuration management
