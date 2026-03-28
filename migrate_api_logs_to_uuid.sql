-- Migration Script: Convert api_logs from SERIAL to UUID
-- Date: 2026-03-28
-- Description: Drop old api_logs table and recreate with UUID

-- Step 1: Drop existing table (WARNING: This will delete all existing logs)
DROP TABLE IF EXISTS api_logs;

-- Step 2: Create new table with UUID
CREATE TABLE api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time BIGINT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    access_from VARCHAR(50),
    user_id UUID,
    company_id UUID,
    branch_id UUID,
    request_body TEXT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Step 3: Create indexes
CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_company_id ON api_logs(company_id);
CREATE INDEX idx_api_logs_branch_id ON api_logs(branch_id);
CREATE INDEX idx_api_logs_access_from ON api_logs(access_from);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at);
CREATE INDEX idx_api_logs_deleted_at ON api_logs(deleted_at);

-- Done!
-- Note: All existing logs have been deleted. New logs will use UUID.
