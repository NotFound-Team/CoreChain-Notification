CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    fcm_token TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP,
    retry_count INT DEFAULT 0,
    task_id VARCHAR(100),
    project_id VARCHAR(100),
    priority INT,
    CONSTRAINT chk_status CHECK (status IN ('pending', 'sent', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_task_id ON notifications(task_id) WHERE task_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_type_status ON notifications(notification_type, status);

CREATE TABLE IF NOT EXISTS user_fcm_tokens (
    user_id VARCHAR(100) PRIMARY KEY,
    fcm_token TEXT NOT NULL,
    platform VARCHAR(20),
    last_updated TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fcm_tokens_updated ON user_fcm_tokens(last_updated DESC);
