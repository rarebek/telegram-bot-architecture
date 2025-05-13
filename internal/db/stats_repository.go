package db

import (
	"context"
	"fmt"
	"time"
)

// UpdateDailyStats updates or creates the statistics for the current day
func (db *DB) UpdateDailyStats(ctx context.Context) error {
	today := time.Now().Format("2006-01-02")

	// Get today's stats from messages and users tables
	query := `
		WITH today_stats AS (
			SELECT
				COUNT(*) as total_messages,
				COUNT(*) FILTER (WHERE is_command = true) as total_commands,
				COUNT(DISTINCT user_id) as active_users
			FROM messages
			WHERE DATE(created_at) = $1::date
		),
		new_users AS (
			SELECT COUNT(*) as count
			FROM users
			WHERE DATE(created_at) = $1::date
		)
		INSERT INTO statistics (date, total_commands, total_messages, active_users, new_users, updated_at)
		SELECT
			$1::date,
			COALESCE(today_stats.total_commands, 0),
			COALESCE(today_stats.total_messages, 0),
			COALESCE(today_stats.active_users, 0),
			COALESCE(new_users.count, 0),
			NOW()
		FROM today_stats, new_users
		ON CONFLICT (date) DO UPDATE
		SET
			total_commands = COALESCE(EXCLUDED.total_commands, 0),
			total_messages = COALESCE(EXCLUDED.total_messages, 0),
			active_users = COALESCE(EXCLUDED.active_users, 0),
			new_users = COALESCE(EXCLUDED.new_users, 0),
			updated_at = NOW()
	`

	_, err := db.conn.Exec(ctx, query, today)
	if err != nil {
		return fmt.Errorf("failed to update daily stats: %w", err)
	}

	return nil
}

// GetDailyStats retrieves statistics for the last n days
func (db *DB) GetDailyStats(ctx context.Context, days int) ([]map[string]interface{}, error) {
	query := `
		SELECT
			date,
			total_commands,
			total_messages,
			active_users,
			new_users
		FROM statistics
		ORDER BY date DESC
		LIMIT $1
	`

	rows, err := db.conn.Query(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var date time.Time
		var commands, messages, activeUsers, newUsers int

		if err := rows.Scan(&date, &commands, &messages, &activeUsers, &newUsers); err != nil {
			return nil, fmt.Errorf("failed to scan stats row: %w", err)
		}

		stat := map[string]interface{}{
			"date":           date.Format("2006-01-02"),
			"total_commands": commands,
			"total_messages": messages,
			"active_users":   activeUsers,
			"new_users":      newUsers,
		}

		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stats rows: %w", err)
	}

	return stats, nil
}
