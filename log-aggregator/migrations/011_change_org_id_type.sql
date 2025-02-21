-- Schema Version: 1.0.0
-- Created: 2025-02-21
-- Description: Change organization_id columns from UUID to VARCHAR(24) for MongoDB compatibility

-- Drop existing foreign key constraints and indexes
ALTER TABLE logs DROP CONSTRAINT IF EXISTS logs_organization_id_fkey;
ALTER TABLE alerts DROP CONSTRAINT IF EXISTS alerts_organization_id_fkey;
ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_organization_id_fkey;

DROP INDEX IF EXISTS idx_logs_org_id;
DROP INDEX IF EXISTS idx_alerts_org_id;
DROP INDEX IF EXISTS idx_api_keys_org_id;
DROP INDEX IF EXISTS idx_logs_org_timestamp;
DROP INDEX IF EXISTS idx_alerts_org_status;

-- Modify column types
ALTER TABLE organizations 
    ALTER COLUMN id TYPE VARCHAR(24);

ALTER TABLE logs 
    ALTER COLUMN organization_id TYPE VARCHAR(24);

ALTER TABLE alerts 
    ALTER COLUMN organization_id TYPE VARCHAR(24);

ALTER TABLE api_keys 
    ALTER COLUMN organization_id TYPE VARCHAR(24);

-- Re-establish foreign key constraints
ALTER TABLE logs
    ADD CONSTRAINT logs_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

ALTER TABLE alerts
    ADD CONSTRAINT alerts_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

ALTER TABLE api_keys
    ADD CONSTRAINT api_keys_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

-- Recreate indexes
CREATE INDEX idx_logs_org_id ON logs(organization_id);
CREATE INDEX idx_alerts_org_id ON alerts(organization_id);
CREATE INDEX idx_api_keys_org_id ON api_keys(organization_id);
CREATE INDEX idx_logs_org_timestamp ON logs(organization_id, timestamp);
CREATE INDEX idx_alerts_org_status ON alerts(organization_id, status);

-- Down migration
-- In case we need to rollback, convert back to UUID
ALTER TABLE logs DROP CONSTRAINT IF EXISTS logs_organization_id_fkey;
ALTER TABLE alerts DROP CONSTRAINT IF EXISTS alerts_organization_id_fkey;
ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_organization_id_fkey;

DROP INDEX IF EXISTS idx_logs_org_id;
DROP INDEX IF EXISTS idx_alerts_org_id;
DROP INDEX IF EXISTS idx_api_keys_org_id;
DROP INDEX IF EXISTS idx_logs_org_timestamp;
DROP INDEX IF EXISTS idx_alerts_org_status;

ALTER TABLE organizations ALTER COLUMN id TYPE UUID USING id::uuid;
ALTER TABLE logs ALTER COLUMN organization_id TYPE UUID USING organization_id::uuid;
ALTER TABLE alerts ALTER COLUMN organization_id TYPE UUID USING organization_id::uuid;
ALTER TABLE api_keys ALTER COLUMN organization_id TYPE UUID USING organization_id::uuid;

ALTER TABLE logs
    ADD CONSTRAINT logs_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

ALTER TABLE alerts
    ADD CONSTRAINT alerts_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

ALTER TABLE api_keys
    ADD CONSTRAINT api_keys_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE;

CREATE INDEX idx_logs_org_id ON logs(organization_id);
CREATE INDEX idx_alerts_org_id ON alerts(organization_id);
CREATE INDEX idx_api_keys_org_id ON api_keys(organization_id);
CREATE INDEX idx_logs_org_timestamp ON logs(organization_id, timestamp);
CREATE INDEX idx_alerts_org_status ON alerts(organization_id, status);
