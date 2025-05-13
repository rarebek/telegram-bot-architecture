package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// DB represents a database connection
type DB struct {
	conn *pgx.Conn
}

// New creates a new database connection
func New(dsn string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{
		conn: conn,
	}, nil
}

// Close closes the database connection
func (db *DB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}

// Ping checks the database connection
func (db *DB) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

// GetConnection returns the underlying database connection
func (db *DB) GetConnection() *pgx.Conn {
	return db.conn
}
