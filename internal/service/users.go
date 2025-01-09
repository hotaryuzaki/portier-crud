package service

import (
	"context"
	"fmt"
	"log"
	"portier/pkg/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	GenderStr string    `json:"gender"` // Temporary field to hold the string value
	Gender    bool      `json:"-"`      // true = male, false = female. This is to make the gender always flexible in the Frontend
	IDNumber  string    `json:"id_number"`
	UserImage string    `json:"user_image"`
	TenantID  int       `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// ConvertGender converts the GenderStr to a boolean value
func (u *User) ConvertGender() error {
	if u.GenderStr == "1" {
		u.Gender = true
	} else if u.GenderStr == "0" {
		u.Gender = false
	} else {
		return fmt.Errorf("invalid gender value")
	}
	return nil
}

// GetAllUsers fetches all users
func GetAllUsers(limit, offset int) ([]User, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username, email, password, name, gender, id_number, user_image, tenant_id, created_at, is_active 
			  FROM users 
			  ORDER BY id 
			  LIMIT $1 OFFSET $2`
	rows, err := dbConn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Gender, &user.IDNumber, &user.UserImage, &user.TenantID, &user.CreatedAt, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID fetches a user by their ID
func GetUserByID(id int) (User, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Add context

	var user User

	query := `SELECT id, username, email, password, name, gender, id_number, user_image, tenant_id, created_at, is_active FROM users WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Gender, &user.IDNumber, &user.UserImage, &user.TenantID, &user.CreatedAt, &user.IsActive)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user
func CreateUser(user User) (User, error) {
	tenants, err := GetAllTenants(1, 0)
	if err != nil {
		fmt.Println("Error getting tenants 000:", err)
		return user, err
	}

	if len(tenants) > 0 {
		tenantID := tenants[0].ID
		user.TenantID = tenantID
	} else {
		fmt.Println("No tenants found")
		return user, err
	}

	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Replace the original password with the hashed one
	user.Password = string(hashedPassword)
	// Explicitly set the default value for IsActive
	user.IsActive = true

	// SQL query to insert a new user
	query := `INSERT INTO users (username, email, password, name, gender, id_number, user_image, tenant_id, created_at, is_active) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	// Insert user data into the database and retrieve the generated ID
	var id int
	err = dbConn.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Name, user.Gender, user.IDNumber, user.UserImage, user.TenantID, time.Now(), true).Scan(&id)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return User{}, fmt.Errorf("failed to create user: %v", err) // Wrap the error with more context
	}

	user.ID = id     // Set the generated user ID
	return user, nil // Return the created user
}

// UpdateUser updates a user's information
func UpdateUser(id int, user User) (User, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Replace the original password with the hashed one
	user.Password = string(hashedPassword)

	query := `UPDATE users SET username=$1, email=$2, password=$3, name=$4, gender=$5, id_number=$6, user_image=$7, tenant_id=$8, is_active=$9 WHERE id=$10`
	_, err = dbConn.Exec(ctx, query, user.Username, user.Email, user.Password, user.Name, user.Gender, user.IDNumber, user.UserImage, user.TenantID, user.IsActive, id)
	if err != nil {
		return User{}, err
	}

	user.ID = id
	return user, nil
}

// DeleteUser deletes a user
func DeleteUser(id int) error {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Context for the query

	query := `DELETE FROM users WHERE id=$1`
	_, err := dbConn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
