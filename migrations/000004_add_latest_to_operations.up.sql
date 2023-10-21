ALTER TABLE operations ADD COLUMN IF NOT EXISTS latest BOOLEAN NOT NULL DEFAULT false;
CREATE UNIQUE INDEX IF NOT EXISTS idx_operations_latest_true ON operations (latest) WHERE latest;
