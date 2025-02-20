-- Schema Version: 1.0.0
-- Created: 2025-02-19
-- Description: Add user tracking to logs and alerts

-- Add user_id column to logs table
ALTER TABLE logs
ADD COLUMN user_id VARCHAR(24);

-- Add user_id column to alerts table
ALTER TABLE alerts
ADD COLUMN user_id VARCHAR(24);

-- Add indexes for efficient querying
CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_alerts_user_id ON alerts(user_id);

-- Add composite indexes for common queries
CREATE INDEX idx_logs_user_timestamp ON logs(user_id, timestamp);
CREATE INDEX idx_alerts_user_status ON alerts(user_id, status);
