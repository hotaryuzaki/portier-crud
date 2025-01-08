package service

import (
	"context"
	"fmt"
	"portier/pkg/db"
	"time"
)

type Key struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by,omitempty"` // Optional, nullable field
	IsActive  bool      `json:"is_active"`
}

// GetAllKeys fetches all keys
func GetAllKeys() ([]Key, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	var keys []Key

	rows, err := dbConn.Query(ctx, "SELECT id, name, created_at, created_by, is_active FROM keys")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key Key
		if err := rows.Scan(&key.ID, &key.Name, &key.CreatedAt, &key.CreatedBy, &key.IsActive); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

// GetKeysByID fetches a key by their ID
func GetKeysByID(id int) (Key, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Add context

	var key Key

	query := `SELECT id, name, created_at, created_by, is_active FROM keys WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&key.ID, &key.Name, &key.CreatedAt, &key.CreatedBy, &key.IsActive)
	if err != nil {
		return Key{}, err
	}

	return key, nil
}

// CreateKey creates a new key
func CreateKey(key Key) (Key, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Explicitly set the default value for IsActive
	key.IsActive = true

	query := `INSERT INTO keys (name, created_at, is_active, created_by) 
						VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := dbConn.QueryRow(ctx, query, key.Name, time.Now(), key.IsActive, key.CreatedBy).Scan(&id)
	if err != nil {
		return Key{}, fmt.Errorf("failed to create key: %v", err)
	}

	key.ID = id
	return key, nil
}

// UpdateKey updates a key's information
func UpdateKey(id int, key Key) (Key, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `UPDATE keys SET name=$1, is_active=$2, created_by=$3 WHERE id=$4`
	_, err := dbConn.Exec(ctx, query, key.Name, key.IsActive, key.CreatedBy, id)
	if err != nil {
		return Key{}, err
	}

	key.ID = id
	return key, nil
}

// DeleteKey deletes a key
func DeleteKey(id int) error {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `DELETE FROM keys WHERE id=$1`
	_, err := dbConn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
