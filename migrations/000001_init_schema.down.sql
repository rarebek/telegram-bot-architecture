-- Drop tables in reverse order
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS statistics;
DROP TABLE IF EXISTS users;

-- Drop indexes (though they would be dropped with the tables)
DROP INDEX IF EXISTS idx_users_telegram_id;
DROP INDEX IF EXISTS idx_messages_user_id;
DROP INDEX IF EXISTS idx_messages_created_at;
