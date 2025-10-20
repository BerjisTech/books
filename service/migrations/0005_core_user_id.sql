-- Add Core user linkage columns without altering existing UUID columns
ALTER TABLE authors ADD COLUMN IF NOT EXISTS core_user_id TEXT;
CREATE INDEX IF NOT EXISTS idx_authors_core_user ON authors(core_user_id);

ALTER TABLE publishers ADD COLUMN IF NOT EXISTS core_user_id TEXT;
CREATE INDEX IF NOT EXISTS idx_publishers_core_user ON publishers(core_user_id);

