package db

import (
	"context"
	"fmt"
	"time"
)

// User represents a user in the database
type User struct {
	ID              int64
	TelegramID      int64
	Username        string
	FirstName       string
	LastName        string
	CommandCount    int
	MessageCount    int
	LastInteraction time.Time
	CreatedAt       time.Time
}

// CreateUser inserts a new user into the database
func (db *DB) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, created_at, last_interaction)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	// Set current time for created_at and last_interaction if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.LastInteraction.IsZero() {
		user.LastInteraction = time.Now()
	}

	return db.conn.QueryRow(
		ctx,
		query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.LastInteraction,
	).Scan(&user.ID)
}

// GetUserByTelegramID retrieves a user by their Telegram ID
func (db *DB) GetUserByTelegramID(ctx context.Context, telegramID int64) (*User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name,
		       command_count, message_count, last_interaction, created_at
		FROM users
		WHERE telegram_id = $1
	`

	user := &User{}
	err := db.conn.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.CommandCount,
		&user.MessageCount,
		&user.LastInteraction,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetOrCreateUser retrieves a user by Telegram ID or creates them if they don't exist
func (db *DB) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*User, error) {
	user, err := db.GetUserByTelegramID(ctx, telegramID)
	if err == nil {
		return user, nil
	}

	// Create new user if not found
	newUser := &User{
		TelegramID:      telegramID,
		Username:        username,
		FirstName:       firstName,
		LastName:        lastName,
		CreatedAt:       time.Now(),
		LastInteraction: time.Now(),
	}

	if err := db.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// UpdateUserActivity updates a user's activity metrics
func (db *DB) UpdateUserActivity(ctx context.Context, telegramID int64, isCommand bool) error {
	var query string
	if isCommand {
		query = `
			UPDATE users
			SET command_count = command_count + 1, last_interaction = $2
			WHERE telegram_id = $1
		`
	} else {
		query = `
			UPDATE users
			SET message_count = message_count + 1, last_interaction = $2
			WHERE telegram_id = $1
		`
	}

	_, err := db.conn.Exec(ctx, query, telegramID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user activity: %w", err)
	}

	return nil
}

// GetTopUsers returns the most active users based on total activity
func (db *DB) GetTopUsers(ctx context.Context, limit int) ([]*User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name,
		       command_count, message_count, last_interaction, created_at
		FROM users
		ORDER BY (command_count + message_count) DESC
		LIMIT $1
	`

	rows, err := db.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.CommandCount,
			&user.MessageCount,
			&user.LastInteraction,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
