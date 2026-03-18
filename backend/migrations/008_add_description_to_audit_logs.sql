-- Add description column to audit_logs table
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS description TEXT;

-- Create index on description for better search performance (optional)
CREATE INDEX IF NOT EXISTS idx_audit_logs_description ON audit_logs(description) WHERE description IS NOT NULL;
