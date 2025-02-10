-- Make organization_id nullable for alerts and logs when multi-tenancy is disabled
ALTER TABLE alerts DROP CONSTRAINT alerts_organization_id_fkey;
ALTER TABLE alerts ALTER COLUMN organization_id DROP NOT NULL;
ALTER TABLE alerts ADD CONSTRAINT alerts_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE logs DROP CONSTRAINT logs_organization_id_fkey;
ALTER TABLE logs ALTER COLUMN organization_id DROP NOT NULL;
ALTER TABLE logs ADD CONSTRAINT logs_organization_id_fkey 
    FOREIGN KEY (organization_id) 
    REFERENCES organizations(id) 
    ON DELETE CASCADE
    DEFERRABLE INITIALLY DEFERRED;

-- Down migration
-- ALTER TABLE alerts DROP CONSTRAINT alerts_organization_id_fkey;
-- ALTER TABLE alerts ALTER COLUMN organization_id SET NOT NULL;
-- ALTER TABLE alerts ADD CONSTRAINT alerts_organization_id_fkey 
--     FOREIGN KEY (organization_id) 
--     REFERENCES organizations(id) 
--     ON DELETE CASCADE;

-- ALTER TABLE logs DROP CONSTRAINT logs_organization_id_fkey;
-- ALTER TABLE logs ALTER COLUMN organization_id SET NOT NULL;
-- ALTER TABLE logs ADD CONSTRAINT logs_organization_id_fkey 
--     FOREIGN KEY (organization_id) 
--     REFERENCES organizations(id) 
--     ON DELETE CASCADE;
