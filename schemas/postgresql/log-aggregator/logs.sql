-- Schema Version: 1.0.0
-- Created: 2025-01-31
-- Last Modified: 2025-01-31
-- Description: Core schema for log aggregation system

-- Logs table schema
CREATE TABLE IF NOT EXISTS logs (
    id VARCHAR(36) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    host VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    level VARCHAR(50) NOT NULL,
    metadata JSONB,
    process_count INTEGER,
    total_cpu_percent FLOAT,
    total_memory_usage BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Process logs table schema
CREATE TABLE IF NOT EXISTS process_logs (
    id VARCHAR(36) PRIMARY KEY,
    log_id VARCHAR(36) REFERENCES logs(id),
    name VARCHAR(255) NOT NULL,
    pid INTEGER NOT NULL,
    cpu_percent FLOAT,
    memory_usage BIGINT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_logs_host ON logs(host);
CREATE INDEX IF NOT EXISTS idx_process_logs_log_id ON process_logs(log_id);
CREATE INDEX IF NOT EXISTS idx_process_logs_name ON process_logs(name);

-- Schema Documentation:
-- 
-- logs table:
-- - id: UUID for the log entry
-- - timestamp: When the log event occurred
-- - host: Source host/system name
-- - message: The actual log message
-- - level: Log level (e.g., INFO, ERROR)
-- - metadata: Additional JSON data
-- - process_count: Number of processes at time of log
-- - total_cpu_percent: Overall CPU usage
-- - total_memory_usage: Overall memory usage in bytes
-- - created_at: Record creation timestamp
--
-- process_logs table:
-- - id: UUID for the process log entry
-- - log_id: Reference to parent log entry
-- - name: Process name
-- - pid: Process ID
-- - cpu_percent: Process CPU usage
-- - memory_usage: Process memory usage in bytes
-- - status: Process status
-- - created_at: Record creation timestamp
