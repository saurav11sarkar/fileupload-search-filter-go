ALTER TABLE users
  DROP COLUMN IF EXISTS reset_code,
  DROP COLUMN IF EXISTS reset_code_expires_at;
