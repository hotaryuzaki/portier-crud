package service

import (
	"context"
	"fmt"
	"log"
	"portier/pkg/db"
	"time"
)

type Copy struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	KeyID     int       `json:"key_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by,omitempty"` // Optional, nullable field
	IsActive  bool      `json:"is_active"`
}

// GetAllCopiesResponse represents the response structure for GetAllCopies
type GetAllCopiesResponse struct {
	Copies     []Copy `json:"copies"`
	TotalPages int    `json:"totalPages"`
}

// GetAllCopies fetches all copies
func GetAllCopies(limit, offset int) (GetAllCopiesResponse, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query to get the total count of copies
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM copies`
	if err := dbConn.QueryRow(ctx, countQuery).Scan(&totalCount); err != nil {
		return GetAllCopiesResponse{}, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit

	// Query to get the paginated copies
	query := `SELECT id, name, key_id, created_at, created_by, is_active 
			  FROM copies 
			  ORDER BY id 
			  LIMIT $1 OFFSET $2`
	rows, err := dbConn.Query(ctx, query, limit, offset)
	if err != nil {
		return GetAllCopiesResponse{}, err
	}
	defer rows.Close()

	var copies []Copy
	for rows.Next() {
		var copy Copy
		if err := rows.Scan(&copy.ID, &copy.Name, &copy.KeyID, &copy.CreatedAt, &copy.CreatedBy, &copy.IsActive); err != nil {
			return GetAllCopiesResponse{}, err
		}
		copies = append(copies, copy)
	}

	if err := rows.Err(); err != nil {
		return GetAllCopiesResponse{}, err
	}

	return GetAllCopiesResponse{
		Copies:     copies,
		TotalPages: totalPages,
	}, nil
}

// GetCopyByID fetches a copy by its ID
func GetCopyByID(id int) (Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Add context

	var copy Copy

	query := `SELECT id, name, key_id, created_at, created_by, is_active FROM copies WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&copy.ID, &copy.Name, &copy.KeyID, &copy.CreatedAt, &copy.CreatedBy, &copy.IsActive)
	if err != nil {
		return Copy{}, err
	}

	return copy, nil
}

// CreateCopy creates a new copy
func CreateCopy(copy Copy) (Copy, error) {
	response, err := GetAllKeys(1, 0)
	if err != nil {
		fmt.Println("Error getting keys:", err)
		return copy, err
	}

	if len(response.Keys) > 0 {
		keyID := response.Keys[0].ID
		copy.KeyID = keyID
	} else {
		fmt.Println("No keys found")
		return copy, err
	}

	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Explicitly set the default value for IsActive
	copy.IsActive = true
	copy.CreatedBy = 1

	query := `INSERT INTO copies (name, key_id, created_at, created_by, is_active) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err = dbConn.QueryRow(ctx, query, copy.Name, copy.KeyID, time.Now(), copy.CreatedBy, copy.IsActive).Scan(&id)
	if err != nil {
		log.Printf("Error creating copy: %v", err)
		return Copy{}, fmt.Errorf("failed to create copy: %v", err)
	}

	copy.ID = id
	return copy, nil
}

// UpdateCopy updates a copy's information
func UpdateCopy(id int, copy Copy) (Copy, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Explicitly set the default value for IsActive
	copy.IsActive = true

	query := `UPDATE copies SET name=$1, is_active=$2 WHERE id=$3`
	_, err := dbConn.Exec(ctx, query, copy.Name, copy.IsActive, id)
	if err != nil {
		return Copy{}, err
	}

	copy.ID = id
	return copy, nil
}

// DeleteCopy deletes a copy
func DeleteCopy(id int) error {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `DELETE FROM copies WHERE id=$1`
	_, err := dbConn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
