-- Schema Version: 1.0.0
-- Created: 2025-02-21
-- Description: Add system organization and protection mechanisms

-- Insert system organization
INSERT INTO organizations (id, name, status, settings) 
VALUES ('system', 'System', 'active', '{"type": "system", "protected": true}'::jsonb)
ON CONFLICT (id) DO NOTHING;

-- Create function to prevent system organization modification
CREATE OR REPLACE FUNCTION prevent_system_org_modification()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.id = 'system' THEN
        RAISE EXCEPTION 'Cannot modify or delete system organization';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to protect system organization
CREATE TRIGGER protect_system_org
BEFORE UPDATE OR DELETE ON organizations
FOR EACH ROW
EXECUTE FUNCTION prevent_system_org_modification();

-- Down migration
DROP TRIGGER IF EXISTS protect_system_org ON organizations;
DROP FUNCTION IF EXISTS prevent_system_org_modification();
DELETE FROM organizations WHERE id = 'system';
