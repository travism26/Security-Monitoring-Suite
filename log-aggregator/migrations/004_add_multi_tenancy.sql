-- Add multi-tenancy support

-- Create organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    settings JSONB DEFAULT '{}'::jsonb
);

-- Create API keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    key_type VARCHAR(50) NOT NULL CHECK (key_type IN ('agent', 'customer')),
    key_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'revoked')),
    permissions JSONB DEFAULT '{}'::jsonb
);

-- Add tenant ID to logs table
ALTER TABLE logs 
ADD COLUMN organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

-- Add tenant ID to alerts table
ALTER TABLE alerts
ADD COLUMN organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

-- Add indexes for tenant-based queries
CREATE INDEX idx_logs_org_id ON logs(organization_id);
CREATE INDEX idx_alerts_org_id ON alerts(organization_id);
CREATE INDEX idx_api_keys_org_id ON api_keys(organization_id);
CREATE UNIQUE INDEX idx_api_keys_hash ON api_keys(key_hash);

-- Add composite indexes for common tenant queries
CREATE INDEX idx_logs_org_timestamp ON logs(organization_id, timestamp);
CREATE INDEX idx_alerts_org_status ON alerts(organization_id, status);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Down migration
-- DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
-- DROP FUNCTION IF EXISTS update_updated_at_column();
-- DROP INDEX IF EXISTS idx_logs_org_timestamp;
-- DROP INDEX IF EXISTS idx_alerts_org_status;
-- DROP INDEX IF EXISTS idx_api_keys_hash;
-- DROP INDEX IF EXISTS idx_api_keys_org_id;
-- DROP INDEX IF EXISTS idx_alerts_org_id;
-- DROP INDEX IF EXISTS idx_logs_org_id;
-- ALTER TABLE alerts DROP COLUMN IF EXISTS organization_id;
-- ALTER TABLE logs DROP COLUMN IF EXISTS organization_id;
-- DROP TABLE IF EXISTS api_keys;
-- DROP TABLE IF EXISTS organizations;
