-- Add user_id column to transactions table
ALTER TABLE transactions ADD COLUMN user_id TEXT NOT NULL REFERENCES users(id);

-- Create an index on user_id for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
