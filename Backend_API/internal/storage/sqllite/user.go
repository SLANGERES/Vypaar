package sqllite

import (
	"database/sql"
	"fmt"

	"log/slog"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/slangeres/Vypaar/backend_API/internal/config"
	"github.com/slangeres/Vypaar/backend_API/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserDbConnection struct {
	db *sql.DB
}

// InitUserDb initializes the user DB with SQLite
func InitUserDb(cnf *config.Config) (*UserDbConnection, error) {
	// Correct quotes around "sqlite3"
	db, err := sql.Open("sqlite3", cnf.UserStoragePath)

	if err != nil {
		slog.Error("Unable to connect to the user database", slog.String("error", err.Error()))
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			shopID TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		slog.Error("Unable to create table in Database", slog.String("error", err.Error()))
		return nil, err
	}

	return &UserDbConnection{db: db}, nil
}

func (dbc *UserDbConnection) Login(email string, password string) (int64, error) {
	var userdata types.User

	stmt, err := dbc.db.Prepare(`SELECT id, name, password FROM user WHERE email=$1`)
	if err != nil {
		return 0, fmt.Errorf("error preparing db: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).Scan(&userdata.Id, &userdata.Name, &userdata.Password)
	if err != nil {
		return 0, fmt.Errorf("invalid credentials: %w", err)
	}

	// Use bcrypt to compare passwords securely
	err = bcrypt.CompareHashAndPassword([]byte(userdata.Password), []byte(password))
	if err != nil {
		return 0, fmt.Errorf("wrong password")
	}

	return int64(userdata.Id), nil
}
func (dbc *UserDbConnection) Signup(name string, email string, password string, shopID string) (int64, error) {
	// Hash the password using bcrypt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Prepare SQL statement

	stmt, err := dbc.db.Prepare(`INSERT INTO user(name, email, password, shopID) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the insert with hashed password
	result, err := stmt.Exec(name, email, string(hashedPassword), shopID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}
