-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    command_count INT DEFAULT 0,
    message_count INT DEFAULT 0,
    last_interaction TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    message_text TEXT,
    is_command BOOLEAN DEFAULT FALSE,
    message_type VARCHAR(50),
    chat_id BIGINT NOT NULL,
    chat_type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create statistics table
CREATE TABLE IF NOT EXISTS statistics (
    id SERIAL PRIMARY KEY,
    date DATE UNIQUE NOT NULL,
    total_commands INT DEFAULT 0,
    total_messages INT DEFAULT 0,
    active_users INT DEFAULT 0,
    new_users INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index for faster telegram_id lookups
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);

-- Create index for faster message searches by user_id
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);

-- Create index for faster message time-based queries
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
