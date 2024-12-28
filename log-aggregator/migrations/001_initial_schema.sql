CREATE TABLE IF NOT EXISTS logs (
    id VARCHAR(36) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    host VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    level VARCHAR(50) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_logs_host ON logs(host);
CREATE INDEX IF NOT EXISTS idx_process_logs_log_id ON process_logs(log_id);
CREATE INDEX IF NOT EXISTS idx_process_logs_name ON process_logs(name); 