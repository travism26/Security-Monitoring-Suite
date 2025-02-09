-- Add API key tracking to logs table

-- Add api_key column to logs table
ALTER TABLE logs
ADD COLUMN api_key VARCHAR(255);

-- Create index for efficient querying by API key
CREATE INDEX idx_logs_api_key ON logs(api_key);

-- Down migration
-- DROP INDEX IF EXISTS idx_logs_api_key;
-- ALTER TABLE logs DROP COLUMN IF EXISTS api_key;
