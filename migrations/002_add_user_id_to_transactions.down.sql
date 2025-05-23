-- Drop the index first
DROP INDEX IF EXISTS idx_transactions_user_id;

-- Remove the user_id column
ALTER TABLE transactions DROP COLUMN IF EXISTS user_id;
