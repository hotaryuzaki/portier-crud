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

// GetAllUsersResponse represents the response structure for GetAllUsers
type GetAllUsersResponse struct {
	Users      []User `json:"users"`
	TotalPages int    `json:"totalPages"`
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

// ConvertGenderToStr converts the Gender boolean to a string value
func (u *User) ConvertGenderToStr() string {
	if u.Gender {
		return "1"
	}
	return "0"
}

// GetAllUsers fetches all users
func GetAllUsers(limit, offset int, name, idNumber string) (GetAllUsersResponse, error) {
	dbConn := db.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set default values for name and idNumber if they are empty
	if name == "" {
		name = "%"
	}
	if idNumber == "" {
		idNumber = "%"
	}

	// Build the query with optional search/filter parameters
	query := `SELECT id, username, email, name, gender, id_number, user_image, tenant_id, created_at, is_active 
						FROM users 
						WHERE name ILIKE $1 AND id_number ILIKE $2 
						ORDER BY id LIMIT $3 OFFSET $4`
	args := []interface{}{"%" + name + "%", "%" + idNumber + "%", limit, offset}

	// Query to get the total count of users with the same filters
	countQuery := `SELECT COUNT(*) FROM users WHERE name ILIKE $1 AND id_number ILIKE $2`
	countArgs := []interface{}{"%" + name + "%", "%" + idNumber + "%"}

	var totalCount int
	if err := dbConn.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
		return GetAllUsersResponse{}, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit

	// Query to get the paginated users
	rows, err := dbConn.Query(ctx, query, args...)
	if err != nil {
		return GetAllUsersResponse{}, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Gender, &user.IDNumber, &user.UserImage, &user.TenantID, &user.CreatedAt, &user.IsActive); err != nil {
			return GetAllUsersResponse{}, err
		}
		user.GenderStr = user.ConvertGenderToStr()
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return GetAllUsersResponse{}, err
	}

	return GetAllUsersResponse{
		Users:      users,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID fetches a user by their ID
func GetUserByID(id int) (User, error) {
	// Get a database connection
	dbConn := db.GetConnection()
	ctx := context.Background() // Add context

	var user User

	query := `SELECT id, username, email, name, gender, id_number, user_image, tenant_id, created_at, is_active FROM users WHERE id=$1`
	err := dbConn.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Gender, &user.IDNumber, &user.UserImage, &user.TenantID, &user.CreatedAt, &user.IsActive)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user
func CreateUser(user User) (User, error) {
	response, err := GetAllTenants(1, 0)
	if err != nil {
		fmt.Println("Error getting tenants:", err)
		return user, err
	}

	if len(response.Tenants) > 0 {
		tenantID := response.Tenants[0].ID
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
func UpdateUser(id int, updatedUser User) (User, error) {
	dbConn := db.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the password is provided
	if updatedUser.Password == "" {
		log.Println("Updating user without password")
		// Update user without changing the password
		updateQuery := `UPDATE users SET username=$1, email=$2, name=$3, gender=$4, id_number=$5, user_image=$6, tenant_id=$7, is_active=$8 WHERE id=$9`
		_, err := dbConn.Exec(ctx, updateQuery, updatedUser.Username, updatedUser.Email, updatedUser.Name, updatedUser.Gender, updatedUser.IDNumber, updatedUser.UserImage, updatedUser.TenantID, updatedUser.IsActive, id)
		if err != nil {
			return User{}, fmt.Errorf("failed to update user: %v", err)
		}
	} else {
		log.Println("Updating user with password")
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("failed to hash password: %v", err)
		}
		updatedUser.Password = string(hashedPassword)

		// Update user with the new password
		updateQuery := `UPDATE users SET username=$1, email=$2, password=$3, name=$4, gender=$5, id_number=$6, user_image=$7, tenant_id=$8, is_active=$9 WHERE id=$10`
		_, err = dbConn.Exec(ctx, updateQuery, updatedUser.Username, updatedUser.Email, updatedUser.Password, updatedUser.Name, updatedUser.Gender, updatedUser.IDNumber, updatedUser.UserImage, updatedUser.TenantID, updatedUser.IsActive, id)
		if err != nil {
			return User{}, fmt.Errorf("failed to update user: %v", err)
		}
	}

	// Return the updated user data
	updatedUser.ID = id
	updatedUser.Password = "" // remove password from the response
	return updatedUser, nil
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
