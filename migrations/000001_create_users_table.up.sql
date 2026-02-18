-- Create users table

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL  DEFAULT NOW()
);

--Create index on email for faster lookups
CREATE INDEX idx_user_email ON users (email);

--Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

--Create trigger to autoupdate updated_at
CREATE TRIGGER update_user_updated_at
   BEFORE UPDATE ON users
   for EACH ROW
    EXECUTE FUNCTION update_updated_at_column();