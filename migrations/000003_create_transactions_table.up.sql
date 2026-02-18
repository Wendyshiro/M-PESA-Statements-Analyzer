-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    receipt_no VARCHAR(50) NOT NULL,
    completion_time TIMESTAMP NOT NULL,
    details TEXT NOT NULL,
    transaction_status VARCHAR(50),
    amount_paid DECIMAL(15, 2),
    amount_withdrawn DECIMAL(15, 2),
    balance DECIMAL(15, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_transactions_job_id ON transactions(job_id);
CREATE INDEX idx_transactions_completion_time ON transactions(completion_time DESC);
CREATE INDEX idx_transactions_category ON transactions(category);
CREATE INDEX idx_transactions_receipt_no ON transactions(receipt_no);