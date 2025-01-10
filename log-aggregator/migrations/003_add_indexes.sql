-- Add indexes for common queries in the logs table

-- Index for ID lookups (though this might already exist as primary key)
CREATE INDEX IF NOT EXISTS idx_logs_id ON logs(id);

-- Index for timestamp ordering (used in List queries)
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);

-- Compound index for timestamp range queries with ordering
CREATE INDEX IF NOT EXISTS idx_logs_timestamp_range ON logs(timestamp DESC) 
WHERE timestamp IS NOT NULL;

-- Index for host-based queries (common filter in monitoring)
CREATE INDEX IF NOT EXISTS idx_logs_host ON logs(host);

-- Index for level-based queries (common filter for error analysis)
CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);

-- Add comment to track migration
COMMENT ON TABLE logs IS 'Added performance indexes for common queries';
