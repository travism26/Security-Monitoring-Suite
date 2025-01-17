# System Monitoring Agent Deployment

## System Requirements

- **Operating Systems**: Linux, Windows, macOS
- **CPU**: Minimum 2 cores
- **Memory**: Minimum 512MB RAM
- **Disk Space**: 100MB available
- **Network**: Stable internet connection
- **Dependencies**:
  - Kafka (optional, for direct metric submission)
  - Docker (for containerized deployment)

## Installation

### Linux Installation

1. Download the latest release:
   ```bash
   wget https://example.com/monitoring-agent/agent-linux-amd64
   ```
2. Make the binary executable:
   ```bash
   chmod +x agent-linux-amd64
   ```
3. Move to system binaries directory:
   ```bash
   sudo mv agent-linux-amd64 /usr/local/bin/monitoring-agent
   ```
4. Create systemd service:

   ```bash
   sudo nano /etc/systemd/system/monitoring-agent.service
   ```

   Add the following content:

   ```ini
   [Unit]
   Description=System Monitoring Agent
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/monitoring-agent
   Restart=always
   User=monitoring
   Group=monitoring

   [Install]
   WantedBy=multi-user.target
   ```

5. Start and enable the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start monitoring-agent
   sudo systemctl enable monitoring-agent
   ```

### Windows Installation

1. Download the Windows executable
2. Run the installer
3. Configure as a Windows Service

### Container Deployment

1. Pull the Docker image:
   ```bash
   docker pull monitoring/agent:latest
   ```
2. Run the container:
   ```bash
   docker run -d \
     --name monitoring-agent \
     -v /path/to/config:/etc/monitoring-agent \
     -v /var/log/monitoring-agent:/var/log/monitoring-agent \
     monitoring/agent:latest
   ```

## Update Process

1. Stop the current agent:
   ```bash
   sudo systemctl stop monitoring-agent
   ```
2. Backup configuration:
   ```bash
   cp /etc/monitoring-agent/config.yaml /etc/monitoring-agent/config.yaml.bak
   ```
3. Download new version:
   ```bash
   wget https://example.com/monitoring-agent/agent-linux-amd64
   ```
4. Replace binary:
   ```bash
   sudo mv agent-linux-amd64 /usr/local/bin/monitoring-agent
   ```
5. Restart service:
   ```bash
   sudo systemctl start monitoring-agent
   ```

## Verification

Check service status:

```bash
sudo systemctl status monitoring-agent
```

View logs:

```bash
journalctl -u monitoring-agent -f
```
