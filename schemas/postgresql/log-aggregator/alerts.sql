-- Schema Version: 1.0.0
-- Created: 2025-01-31
-- Last Modified: 2025-01-31
-- Description: Alert management schema for log aggregation system

-- Alerts table schema
CREATE TABLE IF NOT EXISTS alerts (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    source VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP
);

-- Alert logs relationship table
CREATE TABLE IF NOT EXISTS alert_logs (
    alert_id VARCHAR(36) NOT NULL,
    log_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (alert_id, log_id),
    FOREIGN KEY (alert_id) REFERENCES alerts(id) ON DELETE CASCADE,
    FOREIGN KEY (log_id) REFERENCES logs(id) ON DELETE CASCADE
);

-- Alert metadata table
CREATE TABLE IF NOT EXISTS alert_metadata (
    alert_id VARCHAR(36) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (alert_id, key),
    FOREIGN KEY (alert_id) REFERENCES alerts(id) ON DELETE CASCADE
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_alerts_source ON alerts(source);

-- Schema Documentation:
-- 
-- alerts table:
-- - id: UUID for the alert
-- - title: Alert title/name
-- - description: Detailed alert description
-- - severity: Alert severity level
-- - status: Current alert status
-- - source: Alert source system/component
-- - created_at: Alert creation timestamp
-- - updated_at: Last update timestamp
-- - resolved_at: When the alert was resolved
--
-- alert_logs table:
-- - alert_id: Reference to parent alert
-- - log_id: Reference to associated log entry
--
-- alert_metadata table:
-- - alert_id: Reference to parent alert
-- - key: Metadata key name
-- - value: Metadata value
