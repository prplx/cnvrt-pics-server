ALTER TABLE operations DROP COLUMN IF EXISTS latest;
DROP INDEX IF EXISTS idx_operations_latest_true;
