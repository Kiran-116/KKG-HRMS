-- Add magic link fields to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS magic_token VARCHAR(255),
ADD COLUMN IF NOT EXISTS magic_expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS must_change_password BOOLEAN DEFAULT false;

-- Create index on magic_token for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_magic_token ON users(magic_token) WHERE magic_token IS NOT NULL;
