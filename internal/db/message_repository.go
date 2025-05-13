package db

import (
	"context"
	"fmt"
	"time"
)

// Message represents a message in the database
type Message struct {
	ID          int64
	UserID      int64
	MessageText string
	IsCommand   bool
	MessageType string
	ChatID      int64
	ChatType    string
	CreatedAt   time.Time
}

// SaveMessage saves a message to the database
func (db *DB) SaveMessage(ctx context.Context, message *Message) error {
	query := `
		INSERT INTO messages (user_id, message_text, is_command, message_type, chat_id, chat_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Set current time for created_at if not set
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	return db.conn.QueryRow(
		ctx,
		query,
		message.UserID,
		message.MessageText,
		message.IsCommand,
		message.MessageType,
		message.ChatID,
		message.ChatType,
		message.CreatedAt,
	).Scan(&message.ID)
}

// GetMessagesByUserID retrieves messages sent by a specific user
func (db *DB) GetMessagesByUserID(ctx context.Context, userID int64, limit int) ([]*Message, error) {
	query := `
		SELECT id, user_id, message_text, is_command, message_type, chat_id, chat_type, created_at
		FROM messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := db.conn.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.MessageText,
			&message.IsCommand,
			&message.MessageType,
			&message.ChatID,
			&message.ChatType,
			&message.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}

// GetMessageStats returns statistics about message counts
func (db *DB) GetMessageStats(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT
			COUNT(*) as total_messages,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(*) FILTER (WHERE is_command = true) as command_count,
			COUNT(*) FILTER (WHERE is_command = false) as regular_message_count,
			COUNT(*) FILTER (WHERE chat_type = 'private') as private_messages,
			COUNT(*) FILTER (WHERE chat_type != 'private') as group_messages
		FROM messages
	`

	var totalMessages, uniqueUsers, commandCount, regularMsgCount, privateMsgs, groupMsgs int
	err := db.conn.QueryRow(ctx, query).Scan(
		&totalMessages,
		&uniqueUsers,
		&commandCount,
		&regularMsgCount,
		&privateMsgs,
		&groupMsgs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message stats: %w", err)
	}

	stats := map[string]int{
		"total_messages":    totalMessages,
		"unique_users":      uniqueUsers,
		"command_count":     commandCount,
		"regular_msg_count": regularMsgCount,
		"private_messages":  privateMsgs,
		"group_messages":    groupMsgs,
	}

	return stats, nil
}
