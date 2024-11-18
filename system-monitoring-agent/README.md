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
  "timestamp": 1731674033,
  "timestamp_utc": "2024-11-15T12:33:53Z",
  "host_info": {
    "os": "darwin",
    "arch": "arm64"
  },
  "cpu_usage": 6.48,
  "memory_usage": 10760388608,
  "memory_usage_percent": 55.67,
  "total_memory": 19327352832,
  "disk": {
    "free": 285302497280,
    "total": 494384795648,
    "usage_percent": 42.29,
    "used": 209082298368
  },
  "network": {
    "bytes_received": 57124138206,
    "bytes_received_per_second": 58175.49,
    "bytes_sent": 17783935852,
    "bytes_sent_per_second": 1906.8
  }
}
```

## Viewing Logs

```bash
cat agent.log
```

## Development Status

This project is currently under active development. Upcoming features:

- [ ] Support for Linux systems
- [ ] Process monitoring
- [ ] System temperature monitoring
- [ ] Alert configurations
- [ ] Web interface for metrics visualization
- [ ] Prometheus metrics endpoint
- [ ] Docker containerization

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
