package service

import (
	"context"
	"fmt"
	"log"
	"portier/pkg/db"
	"time"
)

type Tenant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// GetAllTenants fetches all tenants from the database
func GetAllTenants(limit, offset int) ([]Tenant, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, address, status, created_at, is_active
						FROM tenants 
						ORDER BY id 
						LIMIT $1 OFFSET $2`
	rows, err := dbConn.Query(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []Tenant
	for rows.Next() {
		var tenant Tenant
		if err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Address, &tenant.Status, &tenant.CreatedAt, &tenant.IsActive); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

// GetTenantByID fetches a tenant by their ID
func GetTenantByID(id int) (Tenant, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background()

	var tenant Tenant

	query := `SELECT id, name, address, status, created_at, is_active FROM tenants WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&tenant.ID, &tenant.Name, &tenant.Address, &tenant.Status, &tenant.CreatedAt, &tenant.IsActive)
	if err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

// CreateTenant creates a new tenant in the database
func CreateTenant(tenant Tenant) (Tenant, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Explicitly set the default value for IsActive
	tenant.IsActive = true

	query := `INSERT INTO tenants (name, address, status, created_at, is_active) 
						VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err := dbConn.QueryRow(ctx, query, tenant.Name, tenant.Address, tenant.Status, time.Now(), tenant.IsActive).Scan(&id)
	if err != nil {
		log.Printf("Error creating tenant: %v", err)
		return Tenant{}, fmt.Errorf("failed to create tenant: %v", err)
	}

	tenant.ID = id
	return tenant, nil
}

// UpdateTenant updates a tenant in the database
func UpdateTenant(id int, tenant Tenant) (Tenant, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `UPDATE tenants SET name=$1, address=$2, status=$3, is_active=$4 WHERE id=$5`
	_, err := dbConn.Exec(ctx, query, tenant.Name, tenant.Address, tenant.Status, tenant.IsActive, id)
	if err != nil {
		return Tenant{}, err
	}

	tenant.ID = id
	return tenant, nil
}

// DeleteTenant deletes a tenant from the database
func DeleteTenant(id int) error {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `DELETE FROM tenants WHERE id=$1`
	_, err := dbConn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
