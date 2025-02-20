-- Schema Version: 1.0.0
-- Created: 2025-02-18
-- Description: Add views for system metrics aggregation

-- System metrics hourly aggregation view
CREATE OR REPLACE VIEW system_metrics_hourly AS
SELECT 
    date_trunc('hour', timestamp) as time_bucket,
    organization_id,
    host,
    COUNT(*) as sample_count,
    AVG(total_cpu_percent) as avg_cpu_percent,
    MAX(total_cpu_percent) as max_cpu_percent,
    MIN(total_cpu_percent) as min_cpu_percent,
    AVG(total_memory_usage) as avg_memory_usage,
    MAX(total_memory_usage) as max_memory_usage,
    MIN(total_memory_usage) as min_memory_usage,
    AVG(process_count) as avg_process_count,
    MAX(process_count) as max_process_count,
    MIN(process_count) as min_process_count
FROM logs
WHERE total_cpu_percent IS NOT NULL
    AND total_memory_usage IS NOT NULL
    AND process_count IS NOT NULL
GROUP BY 
    date_trunc('hour', timestamp),
    organization_id,
    host;

-- Process metrics hourly aggregation view
CREATE OR REPLACE VIEW process_metrics_hourly AS
SELECT 
    date_trunc('hour', l.timestamp) as time_bucket,
    l.organization_id,
    l.host,
    pl.name as process_name,
    COUNT(*) as sample_count,
    AVG(pl.cpu_percent) as avg_cpu_percent,
    MAX(pl.cpu_percent) as max_cpu_percent,
    MIN(pl.cpu_percent) as min_cpu_percent,
    AVG(pl.memory_usage) as avg_memory_usage,
    MAX(pl.memory_usage) as max_memory_usage,
    MIN(pl.memory_usage) as min_memory_usage,
    array_agg(DISTINCT pl.status) as status_list
FROM logs l
JOIN process_logs pl ON l.id = pl.log_id
WHERE pl.cpu_percent IS NOT NULL
    AND pl.memory_usage IS NOT NULL
GROUP BY 
    date_trunc('hour', l.timestamp),
    l.organization_id,
    l.host,
    pl.name;

-- Add indexes to improve view performance
CREATE INDEX IF NOT EXISTS idx_logs_metrics ON logs (
    timestamp,
    organization_id,
    host,
    total_cpu_percent,
    total_memory_usage,
    process_count
) WHERE 
    total_cpu_percent IS NOT NULL 
    AND total_memory_usage IS NOT NULL 
    AND process_count IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_process_logs_metrics ON process_logs (
    name,
    cpu_percent,
    memory_usage
) WHERE 
    cpu_percent IS NOT NULL 
    AND memory_usage IS NOT NULL;

-- Add function to get latest host metrics
CREATE OR REPLACE FUNCTION get_latest_host_metrics(
    org_id UUID,
    host_filter VARCHAR DEFAULT NULL
)
RETURNS TABLE (
    host VARCHAR,
    last_seen TIMESTAMP,
    total_cpu_percent FLOAT,
    total_memory_usage BIGINT,
    process_count INTEGER
) AS $$
BEGIN
    RETURN QUERY
    WITH latest_logs AS (
        SELECT DISTINCT ON (host) 
            host,
            timestamp as last_seen,
            total_cpu_percent,
            total_memory_usage,
            process_count
        FROM logs
        WHERE organization_id = org_id
            AND (host_filter IS NULL OR host = host_filter)
            AND total_cpu_percent IS NOT NULL
            AND total_memory_usage IS NOT NULL
            AND process_count IS NOT NULL
        ORDER BY host, timestamp DESC
    )
    SELECT * FROM latest_logs;
END;
$$ LANGUAGE plpgsql;
